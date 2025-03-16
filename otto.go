/*
OttO is used to build IoT applications.

# The package provides

  - Drivers for a few different breakout boards meant to run on the
    Raspberry Pi.

  - A higher level device interface that is agnostic to the underlying
    libraries.  Use your favorite gpiocdev, periph.io, tinygo, gobot,
    etc.

  - For exaple GPIO handled by gpiocdev library, new handler for Linux GPIOs.

  - I2C handled by periph.io drivers

  - Serial devices handled by the serial library

  - MQTT messaging amoung IoT stations and control software

  - Messanger (not to be confused with messages) implements a Pub/Sub
    (MQTT or other) interface between components of your application

  - HTTP REST Server for data gathering and configuration

  - Websockets for realtime bidirectional communication with a UI

  - High performance Web server built in to serve interactive UI's
    and modern API's

  - Station manager to manage the stations that make up an entire
    sensor network

  - Data Manager for temporary data caching and interfaces to update
    your favorite cloud based timeseries database

  - Message library for standardized messages built to be communicate
    events and information between pacakges.

  - Security Todo

# Message Based System

The primary communication model for OttO is a messaging system based
on the Pub/Sub model defaulting to MQTT. oTTo is also heavily invested
in HTTP to implement user interfaces and REST/Graph APIs.

Messaging and HTTP use paths to specify the subject of interest. These
paths can be generically reduced to an ordered collection of strings
seperated by slashes '/'.  Both MQTT topics, http URI's and UNIX
filesystems use this same schema which we use the generalize the
identity of the elements we are addressing.

In other words we can generalize the following identities:

For example:

	    File: /home/rusty/data/hb/temperature
		HTTP: /api/data/hb/temperature
		MQTT: ss/station/hb/temperature

The data within the highest level topic temperature can be represented
say by JSON `{ farenhiet: 100.07 }`

### Meta Data (Station Information)

For example when a station comes alive it can provide some information
about itself using the topic:

	```ss/m/be:ef:ca:fe:02/station```

The station will announce itself along with some meta information and
it's capabilities.  The body of the message might look something like
this:

```json

	{
		"id": "be:ef:ca:fe:02",
		"ip": "10.11.24.24",
	    "sensors": [
			"tempc",
			"humidity",
			"light"
		],
		"relays": [
			"heater",
			"light"
		],
	}

```

### Sensor Data

Sensor data takes on the form:

	```ss/d/<station>/<sensor>/<index>```

Where the source is the Station ID publishing the respective data.
The sensor is the type of data being produced (temp, humidity,
lidar, GPS).

The index is optional in situations where they may be more than
one similar device or sensor, for example a couple of rotation
counters on wheels.

The value published by the sensors is typically going to be
floating point, however these values may also be integer or
string values, including nema-0183.

### Control Data

	```ss/c/<source>/<device>/<index>```

This is essentially the same as the sensor except that control
commands are used to have a particular device change, for example
turning a relay on or off.
*/
package otto

import (
	"fmt"
	"log/slog"

	"github.com/sensorstation/otto/data"
	"github.com/sensorstation/otto/messanger"
	"github.com/sensorstation/otto/server"
	"github.com/sensorstation/otto/station"
)

// Controller is a message handler that oversees all interactions
// with the application.
type Controller interface {
	Init()
	Start()
	Stop()
	SetOttO(o *OttO)
	MsgHandler() func(messanger.Msg)
}

// OttO is a large wrapper around the Station, Server,
// DataManager and Messanger, including some convenience functions.
type OttO struct {
	Name string
	*station.Station
	*station.StationManager
	*server.Server
	*data.DataManager
	messanger.Messanger
	Controller

	mock bool
	hub  bool // maybe hub should be a different struct?
	done chan any
}

// global variables and structures
var (
	Version     string
	Interactive bool
)

func init() {
	Version = "0.0.9"
}

func (o *OttO) Mock(m bool) {
	o.mock = m
}

func (o *OttO) Done() chan any {
	return o.done
}

func (o *OttO) SetController(c Controller) {
	o.Controller = c
	c.SetOttO(o)
}

// OttO is a convinience function starting the MQTT and HTTP servers,
// the station manager and other stuff.
func (o *OttO) Init() {
	if o.done != nil {
		// server has already been started
		fmt.Println("Server has already been started")
		return
	}
	o.done = make(chan any)

	// Setup the otto station
	o.Station = station.NewStation(o.Name)
	if o.hub {
		// Allocate and start the station manager
		o.StationManager = station.GetStationManager()
		// Subscribe for station announcements
	}

	o.DataManager = data.NewDataManager()

	// start web server / rest server
	o.Server = server.GetServer()

	if o.Controller != nil {
		o.Controller.Init()
	}
}

func (o *OttO) Start() error {
	go o.Server.Start(o.done)
	if Interactive {
		go o.Stop()
	} else {
		o.Stop()
	}

	if o.StationManager != nil {
		o.StationManager.Start()
	}

	<-o.done
	o.Stop()
	return nil
}

func (o *OttO) Stop() {
	<-o.done
	slog.Info("Done, cleaning up()")

	if err := server.GetServer().Close(); err != nil {
		slog.Error("Failed to close server", "error", err)
	}

	o.Messanger.Close()
}
