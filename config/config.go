package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cpo/hue-alarm/log"
)

type Config struct {
	NotificationURL          string         `json:"notificationURL"`
	SchedulePart             string         `json:"triggerOnSchedulePart"`
	TestMode                 bool           `json:"testMode"`
	UserName                 string         `json:"userName"`
	StatusMessages           StatusMessages `json:"statusMessages"`
	NotificationDelaySeconds int            `json:"notificationDelaySeconds"`
}

type StatusMessages struct {
	Enabled    bool   `json:"enabled"`
	SMTPServer string `json:"smtpServer"`
	SMTPPort   int    `json:"smtpPort"`
	Password   string `json:"password"`
	From       string `json:"from"`
	To         string `json:"to"`
}

type State struct {
	LastArmed    bool       `json:"armed"`
	LastNotified *time.Time `json:"lastNotified,omitifempty"`
}

func (s *State) String() string {
	return fmt.Sprintf("Armed: %t", s.LastArmed)
}

// ReadConfig reads a configuration file.
func ReadConfig(name string, settings interface{}, create bool) {
	log.Info.Println("Read config")
	settingsStr, err := ioutil.ReadFile(name)
	if err != nil {
		if create {
			return
		}
		panic(err)
	}

	err = json.Unmarshal(settingsStr, &settings)
	if err != nil {
		panic(err)
	}
}

// WriteConfig writes a configuration file.
func WriteConfig(name string, settings interface{}, create bool) {
	log.Info.Println("Write config")
	settingsStr, err := json.MarshalIndent(settings, "", "\t")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(name, settingsStr, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
