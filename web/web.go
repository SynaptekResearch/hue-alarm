package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cpo/hue-alarm/alarmmonitor"
	"github.com/cpo/hue-alarm/config"
	"github.com/cpo/hue-alarm/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type webInterface struct {
	*alarmmonitor.AlarmMonitor
}

func (w *webInterface) getConfig(c echo.Context) error {
	log.Debug.Printf("API: GET config\n")
	return c.JSON(200, w.AlarmMonitor.Config)
}

func (w *webInterface) postConfig(c echo.Context) error {
	log.Debug.Printf("API: POST config\n")
	c.Bind(&w.AlarmMonitor.Config)
	// config.WriteConfig()
	config.WriteConfig("settings.json", w.AlarmMonitor.Config, false)
	w.AlarmMonitor.Reload()
	return c.String(200, "Settings saved and active")
}

func (w *webInterface) getStatus(c echo.Context) error {
	log.Debug.Printf("API: GET status\n")
	var status struct {
		Running bool         `json:"running"`
		Status  config.State `json:"status"`
	}
	status.Running = w.AlarmMonitor.Running
	status.Status = w.AlarmMonitor.Status
	return c.JSON(200, status)
}

func (w *webInterface) postTestNotification(c echo.Context) error {
	log.Debug.Printf("API: POST testNotification\n")
	var request struct {
		URL string
	}
	err := c.Bind(&request)
	if err != nil {
		log.Debug.Printf("%s\n", err)
		return c.String(400, "Error sending request to "+request.URL)
	}

	getURL := request.URL
	if strings.Contains(getURL, "%s") {
		getURL = fmt.Sprintf(getURL, url.QueryEscape("Test notification"))
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
	return c.String(200, "Notification request sent to "+request.URL)
}

// Start starts the web interface. Duh.
func Start(monitor *alarmmonitor.AlarmMonitor) {
	log.Info.Printf("Initializing API\n")
	webiface := webInterface{monitor}

	e := echo.New()
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "",
	}))
	e.Use(middleware.BasicAuth(func(username, password string) bool {
		if username == monitor.Config.AdminUserName && password == monitor.Config.AdminPassword {
			return true
		}
		return false
	}))
	e.Static("/", "static")
	e.Static("/modules", "node_modules")
	e.GET("/api/config", webiface.getConfig)
	e.POST("/api/config", webiface.postConfig)
	e.POST("/api/test-notify", webiface.postTestNotification)
	e.GET("/api/status", webiface.getStatus)
	e.Start(":8080")
}
