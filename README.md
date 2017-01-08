# Short description
A simple SMS based alarm using HUE sensors.

# Long description
This project is a simple Philips HUE based alarm system. Whenever a HUE motion sensor is triggered, the application sends a GET request to a predefined 
http address. This can either be a SMS gateway or a HTTP push notification service.

# Support and liability

I an willing to give support, but in a limited way. Issues with compilation / go / glide are specifically NOT supported. I presume you know what you are doing.
Also: this package is delivered without guarantees; that means I am not liable for anything. So if your house burns down or your wife gets pregnant or 
anything related to the use of this software: it's on you.

# SMS providers

Tested with (please provide additional tested providers):

https://www.cmtelecom.com/

# Setup guide

## Compile the binary

A standard git clone of this repository with go and glide installed (google that!) will build an ARM binary by default.

## Install the binary

I use a Beaglebone black to run this program. The platform is cheap, energy efficient and uses Linux.

It should also be possible to use another tiny linux platform like Raspberry PI.

Use the `scp` command to install the binary to the beaglebone and ssh into the beaglebone.

## Create an initial Configuration

Using VI or NANO, create an initial configuration file. See `Configuration` for all settings.

## Create cron entry

Using `crontab -e`, add an entry:

```
* *    *   *   *     /hue-alarm/hue-alarm --config /hue-alarm/settings.json --delay 5 --runs 12 > /hue-alarm/alarm.log
```

## Create an alarm from HUE app

Open the Philips HUE app and go to "Routines" -> "My routines" and add a routine with the configured trigger word in it. See `Configuration -> triggerOnSchedulePart`.

If the routine containing the trigger word is enabled, the cron job will trigger the configured HTTP endpoint when any sensor returns movement.

# Commandline options

```
Usage of ./hue-alarm:
  -config string
    	-config settings.json (default "settings.json")
  -delay int
    	-delay 5 (default 5)
  -dumpsensordata
    	-dumpsensordata
  -init
    	-init
  -runs int
    	-runs 1 (default 1)
```

## --config

Sets the configuration file name. Default is `settings.json`.

## --dumpsensordata

Dumps the sensor data (for debugging purposes).

## --init

Connect to the HUE bridge and obtain a username. This username should be set in the `configuration file`.

## --runs

The number of times to query the sensors

## --delay

The number of seconds to delay before running again. Try not to oversaturate the HUE bridge by setting this too low. The bridge may become unresponsive or 
slow for light button events. 

Also, try to balance the runs * delay to match your crontab entry.

For example: 

If cron runs every 60 seconds, setting runs to 6 and delay to 10 results in 6 requests per minute.


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

