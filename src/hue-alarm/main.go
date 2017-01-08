package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"strings"

	"net/http"

	"net/url"

	"flag"

	"github.com/cpo/go-hue/configuration"
	"github.com/cpo/go-hue/portal"
	"github.com/cpo/go-hue/schedules"
	"github.com/cpo/go-hue/sensors"
)

type config struct {
	NotificationURL string `json:"notificationURL"`
	SchedulePart    string `json:"triggerOnSchedulePart"`
	TestMode        bool   `json:"testMode"`
	UserName        string `json:"userName"`
}

func initializeUser() {
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

func main() {

	configName := "settings.json"
	dumpSensors := false
	initMode := false
	runs := 1
	delay := 0

	flag.StringVar(&configName, "config", "settings.json", "-config settings.json")
	flag.BoolVar(&dumpSensors, "dumpsensordata", false, "-dumpsensordata")
	flag.BoolVar(&initMode, "init", false, "-init")
	flag.IntVar(&runs, "runs", 1, "-runs 1")
	flag.IntVar(&delay, "delay", 5, "-delay 5")
	flag.Parse()

	if initMode {
		initializeUser()
		return
	}

	settingsStr, err := ioutil.ReadFile(configName)
	if err != nil {
		panic(err)
	}

	settings := config{}
	err = json.Unmarshal(settingsStr, &settings)
	if err != nil {
		panic(err)
	}

	settingsStr, _ = json.Marshal(settings)
	fmt.Printf("%s\n", settingsStr)

	p, err := portal.GetPortal()
	if err != nil {
		panic(err)
	}

	r := schedules.New(p[0].InternalIPAddress, settings.UserName)
	if err != nil {
		panic(err)
	}
	allSchedules, err := r.GetAllSchedules()
	if err != nil {
		panic(err)
	}

	alarmEnabled := false
	for i := range allSchedules {
		if strings.Contains(allSchedules[i].Name, settings.SchedulePart) && allSchedules[i].Status == "enabled" {
			alarmEnabled = true
			// TODO: watch all sensors here
		}
	}

	snsrs := sensors.New(p[0].InternalIPAddress, settings.UserName)

	fmt.Printf("ALARM enabled: %t\n", alarmEnabled)
	tripped := false

	for run := 0; alarmEnabled && run < runs; run++ {
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
			if dumpSensors {
				fmt.Printf("Sensor %d %s model %s\n============================\n%s\n", s[si].ID, s[si].Name, s[si].Type, s[si].String())
			}
		}

		fmt.Printf("ALARM trigger: %t\n", alarmTrigger)

		if settings.TestMode || (alarmEnabled && alarmTrigger) {
			getURL := settings.NotificationURL
			if strings.Contains(getURL, "%s") {
				getURL = fmt.Sprintf(settings.NotificationURL, url.QueryEscape(alarmSource))
			}
			if !settings.TestMode && !tripped {
				req, err := http.NewRequest("GET", getURL, nil)
				if err != nil {
					panic(err)
				}
				client := http.Client{}
				response, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				tripped = true
				defer response.Body.Close()
			} else {
				fmt.Printf("TEST MODE enabled: %s\n", getURL)
			}
		}

		time.Sleep(time.Second * time.Duration(delay))
	}

}
