/*
OttO is used to build IoT applications.

# The package provides

  - MQTT messaging amoung IoT stations and control software
  - HTTP REST Server for data gathering and configuration
  - Websockets for realtime bidirectional communication with a UI
  - High performance Web server built in to serve interactive UI's
    and modern API's
  - Station manager to manage the stations that make up an entire
    sensor network
  - Data Manager for temporary data caching and interfaces to update
    your favorite timeseries database
  - A higher level device interface that is agnostic to the underlying
    libraries.  Use your favorite gpiocdev, periph.io, tinygo, gobot,
    etc.
  - Message library for standardized messages built to be communicate
    events and information between pacakges.
  - Messanger (not to be confused with messages) implements a Pub/Sub
    (MQTT or other) interface between components of your application
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
		MQTT: ss/hb/temperature

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

	```ss/d/<source>/<sensor>/<index>```

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

	"github.com/sensorstation/otto/messanger"
	"github.com/sensorstation/otto/server"
	"github.com/sensorstation/otto/station"
)

// global variables and structures
var (
	Done        chan bool
	StationName string
	Version     string
	Interactive bool
)

func init() {
	StationName = "station"
	Version = "0.0.8"
}

func Cleanup() {
	<-Done
	slog.Info("Done, cleaning up()")

	messanger.GetMQTT().Disconnect(1000)
	server.GetServer().Close()
}

// OttO is a convinience function starting the MQTT and HTTP servers,
// the station manager and other stuff.
func OttO() {
	if Done != nil {
		// server has already been started
		fmt.Println("Server has already been started")
		return
	}
	Done = make(chan bool)

	// Allocate and start the station manager
	stations := station.GetStationManager()
	stations.Start()

	// start web server / rest server
	server := server.GetServer()
	go server.Start()
	if Interactive {
		println("go cleanup")
		go Cleanup()
	} else {
		Cleanup()
	}
}
