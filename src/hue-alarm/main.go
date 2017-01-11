package main

import (
	"hue-alarm/alarmmonitor"
)

func main() {
	monitor := alarmmonitor.New()
	monitor.Run()
}
