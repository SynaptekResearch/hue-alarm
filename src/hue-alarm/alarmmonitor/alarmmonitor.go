package alarmmonitor

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"net/http"
	"net/smtp"
	"net/url"

	"hue-alarm/config"
	"strconv"

	"github.com/cpo/go-hue/configuration"
	"github.com/cpo/go-hue/portal"
	"github.com/cpo/go-hue/schedules"
	"github.com/cpo/go-hue/sensors"
)

// AlarmMonitor holds all the fields and methods for the HUE alarm.
type AlarmMonitor struct {
	configName  string
	dumpSensors bool
	initMode    bool
	runs        int
	delay       int
	config      config.Config

	sts     config.State
	tripped bool
}

// New constructs a new AlarmMonitor.
func New() *AlarmMonitor {
	monitor := AlarmMonitor{}
	monitor.initialize()
	return &monitor
}

// initialize initializes all members
func (m *AlarmMonitor) initialize() {
	m.configName = "m.config.json"
	m.dumpSensors = false
	m.initMode = false
	m.runs = 1
	m.delay = 0

	flag.StringVar(&m.configName, "config", "settings.json", "-config settings.json")
	flag.BoolVar(&m.dumpSensors, "dumpsensordata", false, "-dumpsensordata")
	flag.BoolVar(&m.initMode, "init", false, "-init")
	flag.IntVar(&m.runs, "runs", 1, "-runs 1")
	flag.IntVar(&m.delay, "delay", 5, "-delay 5")
	flag.Parse()

	if m.initMode {
		m.initializeUser()
		return
	}

	m.config = config.Config{}
	config.ReadConfig(m.configName, &m.config, false)
	m.sts = config.State{}
	config.ReadConfig("state.json", &m.sts, true)
}

// Initialize HUE bridge to create a user.
func (m *AlarmMonitor) initializeUser() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please press the LINK button on the hue bridge")
	reader.ReadLine()
	p, err := portal.GetPortal()
	if err != nil {
		panic(err)
	}
	c := configuration.New(p[0].InternalIPAddress)
	response, err := c.CreateUser("HueAlarm", "Monitor")
	if err != nil {
		panic(err)
	}
	username := response[0].Success["username"].(string)
	fmt.Printf("Your username is %s\n", username)
}

// Send notification to configured address
func (m *AlarmMonitor) notify(status string) {
	if !m.config.StatusMessages.Enabled {
		return
	}
	msg := "From: " + m.config.StatusMessages.From + "\n" +
		"To: " + m.config.StatusMessages.To + "\n" +
		"Subject: " + status + "\n\n" + status

	err := smtp.SendMail(m.config.StatusMessages.SMTPServer+":"+strconv.Itoa(m.config.StatusMessages.SMTPPort),
		smtp.PlainAuth("", m.config.StatusMessages.From, m.config.StatusMessages.Password, m.config.StatusMessages.SMTPServer),
		m.config.StatusMessages.From, []string{m.config.StatusMessages.To}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

// Send an alert for a real alarm (or fall back to a normal notify in test mode).
func (m *AlarmMonitor) notifyAlarm(source string) {
	onlyNotifyAfter := time.Now().Add(-1 * time.Second)
	if m.sts.LastNotified != nil {
		onlyNotifyAfter = m.sts.LastNotified.Add(time.Duration(m.config.NotificationDelaySeconds) * time.Second)
	}

	timeoutPassed := time.Now().After(onlyNotifyAfter)
	if !m.tripped && timeoutPassed {
		fmt.Printf("Alarm tripped, sending notification.\n")
		if m.config.TestMode {
			fmt.Printf("Test mode, using e-mail.\n")
			m.notify(source)
		} else {
			getURL := m.config.NotificationURL
			if strings.Contains(getURL, "%s") {
				getURL = fmt.Sprintf(m.config.NotificationURL, url.QueryEscape(source))
			}
			req, err := http.NewRequest("GET", getURL, nil)
			if err != nil {
				panic(err)
			}
			client := http.Client{}
			response, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer response.Body.Close()
		}
		now := time.Now()
		m.sts.LastNotified = &now
		m.tripped = true
	} else {
		fmt.Printf("Alarm tripped, NOT sending notification because tripped %t and/or timeoutPassed %t.\n", m.tripped, timeoutPassed)
	}

}

// Run the main loop.
func (m *AlarmMonitor) Run() {

	if m.config.TestMode {
		fmt.Printf("Running in test mode.\n")
	}

	p, err := portal.GetPortal()
	if err != nil {
		panic(err)
	}

	if len(p) < 1 {
		panic("Error contacting hue.")
	}

	// TODO: iterate over multiple HUE brigdes, multiple usernames etc
	r := schedules.New(p[0].InternalIPAddress, m.config.UserName)
	if err != nil {
		panic(err)
	}
	allSchedules, err := r.GetAllSchedules()
	if err != nil {
		panic(err)
	}

	alarmEnabled := false
	for i := range allSchedules {
		if strings.Contains(allSchedules[i].Name, m.config.SchedulePart) && allSchedules[i].Status == "enabled" {
			alarmEnabled = true
		}
	}

	// state change, notify via e-mail
	if m.sts.LastArmed != alarmEnabled {
		enabledStr := "Disabled"
		if alarmEnabled {
			enabledStr = "Enabled"
		}
		notifyStr := "Alarm is now " + enabledStr
		fmt.Printf("State change,notify '%s'\n", enabledStr)
		m.notify(notifyStr)
	}

	// process state
	m.sts.LastArmed = alarmEnabled
	snsrs := sensors.New(p[0].InternalIPAddress, m.config.UserName)

	fmt.Printf("Alarm enabled: %t\n", alarmEnabled)
	m.tripped = false

	for run := 0; alarmEnabled && run < m.runs; run++ {
		s, err := snsrs.GetAllSensors()
		if err != nil {
			panic(err)
		}

		alarmTrigger := false
		alarmSource := ""
		for si := range s {
			if s[si].Type == "ZLLPresence" && s[si].State.Presence {
				alarmTrigger = true
				alarmSource = s[si].Name
			}
			if m.dumpSensors {
				fmt.Printf("Sensor %d %s model %s\n============================\n%s\n", s[si].ID, s[si].Name, s[si].Type, s[si].String())
			}
		}

		fmt.Printf("ALARM trigger: %t\n", alarmTrigger)

		if alarmEnabled && alarmTrigger {
			m.notifyAlarm(alarmSource)
		}

		time.Sleep(time.Second * time.Duration(m.delay))
	}

	config.WriteConfig("state.json", m.sts, true)

}
