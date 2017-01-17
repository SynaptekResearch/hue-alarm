package alarmmonitor

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"net/http"
	"net/smtp"
	"net/url"

	"strconv"

	"github.com/cpo/hue-alarm/config"

	"github.com/cpo/go-hue/configuration"
	"github.com/cpo/go-hue/portal"
	"github.com/cpo/go-hue/schedules"
	"github.com/cpo/go-hue/sensors"
	"github.com/cpo/hue-alarm/log"
)

// AlarmMonitor holds all the fields and methods for the HUE alarm.
type AlarmMonitor struct {
	configName  string
	dumpSensors bool
	initMode    bool
	runs        int
	delay       int
	Config      config.Config

	Running      bool
	Status       config.State
	tripped      bool
	alarmEnabled bool
	reload       bool
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
	m.runs = 4
	m.delay = 5

	flag.StringVar(&m.configName, "config", "settings.json", "-config settings.json")
	flag.BoolVar(&m.dumpSensors, "dumpsensordata", false, "-dumpsensordata")
	flag.BoolVar(&m.initMode, "init", false, "-init")
	flag.IntVar(&m.runs, "runs", 4, "-runs 1")
	flag.IntVar(&m.delay, "delay", 5, "-delay 5")
	flag.Parse()

	if m.initMode {
		m.initializeUser()
		return
	}

	m.Config = config.Config{}
	config.ReadConfig(m.configName, &m.Config, false)
	m.Status = config.State{}
	config.ReadConfig("state.json", &m.Status, true)
}

func (m *AlarmMonitor) Reload() {
	fmt.Println("Reload requested...")
	m.reload = true
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
	log.Info.Printf("Your username is %s\n", username)
}

// Send notification to configured address
func (m *AlarmMonitor) notify(status string) {
	if !m.Config.StatusMessages.Enabled {
		return
	}
	msg := "From: " + m.Config.StatusMessages.From + "\n" +
		"To: " + m.Config.StatusMessages.To + "\n" +
		"Subject: " + status + "\n\n" + status

	err := smtp.SendMail(m.Config.StatusMessages.SMTPServer+":"+strconv.Itoa(m.Config.StatusMessages.SMTPPort),
		smtp.PlainAuth("", m.Config.StatusMessages.From, m.Config.StatusMessages.Password, m.Config.StatusMessages.SMTPServer),
		m.Config.StatusMessages.From, []string{m.Config.StatusMessages.To}, []byte(msg))

	if err != nil {
		log.Info.Printf("smtp error: %s", err)
		return
	}
}

// Send an alert for a real alarm (or fall back to a normal notify in test mode).
func (m *AlarmMonitor) notifyAlarm(source string) {
	onlyNotifyAfter := time.Now().Add(-1 * time.Second)
	if m.Status.LastNotified != nil {
		onlyNotifyAfter = m.Status.LastNotified.Add(time.Duration(m.Config.NotificationDelaySeconds) * time.Second)
	}

	timeoutPassed := time.Now().After(onlyNotifyAfter)
	if !m.tripped && timeoutPassed {
		log.Info.Printf("Alarm tripped, sending notification.\n")
		if m.Config.TestMode {
			log.Info.Printf("Test mode, using e-mail.\n")
			m.notify(source)
		} else {
			getURL := m.Config.NotificationURL
			if strings.Contains(getURL, "%s") {
				getURL = fmt.Sprintf(m.Config.NotificationURL, url.QueryEscape(source))
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
		m.Status.LastNotified = &now
		config.WriteConfig("state.json", m.Status, true)
		m.tripped = true
	} else {
		log.Info.Printf("Alarm tripped, NOT sending notification because tripped %t and/or timeoutPassed %t.\n", m.tripped, timeoutPassed)
	}

}

func (m *AlarmMonitor) Close() {
	m.Running = false
}

// Run the main loop.
func (m *AlarmMonitor) Run() {
	m.Running = true
	defer m.Close()
	for {
		log.Info.Printf("New run [testMode: %t]\n", m.Config.TestMode)

		p, err := portal.GetPortal()
		if err != nil {
			log.Info.Printf("Error: %s\n", err)
			continue
		}

		if len(p) < 1 {
			panic("Error contacting hue.")
		}

		// TODO: iterate over multiple HUE brigdes, multiple usernames etc
		r := schedules.New(p[0].InternalIPAddress, m.Config.UserName)
		if err != nil {
			panic(err)
		}
		allSchedules, err := r.GetAllSchedules()
		if err != nil {
			panic(err)
		}

		m.alarmEnabled = false
		for i := range allSchedules {
			if strings.Contains(allSchedules[i].Name, m.Config.SchedulePart) && allSchedules[i].Status == "enabled" {
				m.alarmEnabled = true
			}
		}

		// state change, notify via e-mail
		if m.Status.LastArmed != m.alarmEnabled {
			enabledStr := "Disabled"
			if m.alarmEnabled {
				enabledStr = "Enabled"
			}
			notifyStr := fmt.Sprintf("Alarm is now %s", enabledStr)
			if m.Config.TestMode {
				notifyStr += ", running in test mode."
			}
			log.Info.Printf("State change, notifying '%s'\n", enabledStr)
			m.notify(notifyStr)
			m.Status.LastArmed = m.alarmEnabled
			config.WriteConfig("state.json", m.Status, true)
		}

		// process state
		snsrs := sensors.New(p[0].InternalIPAddress, m.Config.UserName)

		log.Info.Printf("Alarm enabled %t\n", m.alarmEnabled)
		m.tripped = false

		for run := 0; m.alarmEnabled && run < m.runs && !m.reload; run++ {
			s, err := snsrs.GetAllSensors()
			if err != nil {
				log.Info.Printf("Error: %s\n", err)
				continue
			}

			alarmTrigger := false
			alarmSource := ""
			for si := range s {
				if s[si].Type == "ZLLPresence" && s[si].State.Presence {
					alarmTrigger = true
					alarmSource = s[si].Name
				}
				if m.dumpSensors {
					log.Info.Printf("Sensor %d %s model %s\n============================\n%s\n", s[si].ID, s[si].Name, s[si].Type, s[si].String())
				}
			}

			log.Info.Printf("ALARM trigger %t\n", alarmTrigger)

			if m.alarmEnabled && alarmTrigger {
				m.notifyAlarm(alarmSource)
			}

			time.Sleep(time.Second * time.Duration(m.delay))
		}

		if m.reload {
			log.Info.Printf("Reload finished\n")
			m.reload = false

		}
		if !m.alarmEnabled {
			time.Sleep(time.Second * time.Duration(m.delay))
		}

	}
}
