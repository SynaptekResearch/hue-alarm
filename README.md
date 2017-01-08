# Short description
A simple SMS based alarm using HUE sensors.

# Long description
This project is a simple Philips HUE based alarm system. Whenever a HUE motion sensor is triggered, the application sends a GET request to a predefined 
http address. This can either be a SMS gateway or a HTTP push notification service.

# Configuration
```
{
  "notificationURL": "https://sgw01.cm.nl/gateway.ashx?producttoken=XXX-EDITED-OUT-XXX&body=%s&to=XXX-YOUR-PHONE-NUMBER-XXX&from=HUE&reference=HUE",
  "triggerOnSchedulePart": "(ALARM)",
  "testMode": false,
  "userName": "XXX"
}
```

## notificationURL

The URL to use to notify someone. If the URL contains %s, this placeholder will be replaced with the URL encoded name of the sensor that was triggered.

## triggerOnSchedulePart

If a schedule is active with this string as a part of the name, the GET will be performed (see notificationURL) when a sensor is triggered. You can 
easily setup a custom routine with this string in the name. If the schedule is active, the alarm is armed.

## testMode

Just print the URL used, regardless of active schedules and triggered sensors. 

## userName

The username which is known in your HUE bridge. Run the program using the `--init` parameter to obtain a username.

