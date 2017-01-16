package main

import (
	"github.com/cpo/hue-alarm/alarmmonitor"
	"github.com/cpo/hue-alarm/web"
)

func main() {
	monitor := alarmmonitor.New()
	go monitor.Run()
	web.Start(monitor)
}
