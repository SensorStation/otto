/*
The iote package was designed to help build IoT edge software for on
premise edge devices managing a large number of IoT Stations.

# The package provides

- MQTT messaging amoung IoT stations and control software
- HTTP REST Server for data gathering and configuration
- Websockets for realtime bidirectional communication with UI
- Web server for mondern web based User Interface
- Station manager to track a variety of IoT stations

# Messaging Based

The primary communication model for IoTe is a messaging system based
on MQTT. These messages can be broke into the following categories

- Meta Data that help describe the "network" infrastructure
- Sensor Data that produces data gathered by sensors
- Control Data that represents actions to be performed by stations

The topic format used by MQTT is flexible but generally follows the
following formats:

## ss/m/<source>/<type> -> { station-informaiton }

Where ss/m == sensor station, <source> is the station Id or source
of the message and type represents the specific type of information.

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

import "fmt"

// global variables and structures
var (
	mqtt     *MQTT
	server   *Server
	stations *StationManager
	data     *DataManager
	blasters *MQTTBlasters
	config   *Configuration
	l        *Logger

	Done chan bool
)

func init() {
	config = &Configuration{
		Addr:        ":8011",
		Broker:      "localhost",
		Interactive: true,
	}
}

func GetConfig() *Configuration {
	return config
}

func GetMQTT() *MQTT {
	if mqtt == nil {
		mqtt = NewMQTT()
	}
	return mqtt
}

func GetMQTTBlasters() *MQTTBlasters {
	return blasters
}

func GetDataManager() *DataManager {
	if data == nil {
		data = NewDataManager()
	}
	return data
}

func GetStationManager() *StationManager {
	if stations == nil {
		stations = NewStationManager()
	}
	return stations
}

func GetServer() *Server {
	if server == nil {
		server = NewServer()
	}
	return server
}

func GetLogger() *Logger {
	return l
}

func Cleanup() {

	<-Done
	l.Info("Done, cleaning up()")

	if blasters != nil && blasters.Running {
		blasters.Stop()
	}

	if mqtt != nil {
		mqtt.Disconnect(1000)
	}

	if server != nil {
		server.Close()
	}

	if l != nil {
	}
}

func OttO() {
	if Done != nil {
		// server has already been started
		fmt.Println("Server has already been started")
		return
	}
	Done = make(chan bool)

	// Allocate and start the station manager
	stations := GetStationManager()
	stations.Start()

	mqtt := GetMQTT()
	err := mqtt.Connect()
	if err != nil {
		l.Error("MQTT Failed to connect to broker ", "broker", config.Broker)
	} else {
		mqtt.Subscribe("ss/d/+/+", GetDataManager())
	}

	// start web server / rest server
	server := GetServer()
	go server.Start()
	if config.Interactive {
		println("go cleanup")
		go Cleanup()
	} else {
		Cleanup()
	}
}
