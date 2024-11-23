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

// global variables and structures
var (
	mqtt     *MQTT
	server   *Server
	stations *StationManager
	sensors  *SensorManager
)

func GetMQTT() *MQTT {
	if mqtt == nil {
		mqtt = NewMQTT()
	}
	return mqtt
}

func GetSensorManager() *SensorManager {
	if sensors == nil {
		sensors = NewSensorManager()
	}
	return sensors
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

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"plugin"

// 	gomqtt "github.com/eclipse/paho.mqtt.golang"
// )

// type OttO struct {
// 	// *Dispatcher
// 	// *MQTT
// 	// *Server
// 	// *Store

// 	Plugins []string

// 	Done chan bool
// }

// type controller interface {
// 	Init() error
// }

// var (
// 	O *OttO = nil
// )

// func NewOttO() *OttO {
// 	O = &OttO{
// 		MQTT:   &MQTT{},
// 		Server: &Server{},
// 		Store:  &Store{},
// 	}
// 	return O
// }

// func (o *OttO) Start() {
// 	fmt.Println("TODO remove otto.start")
// }

// func (o *OttO) LoadPlugins(plugins []string) {
// 	// Register some callbacks
// 	// Start the HTTP Server
// 	for _, p := range plugins {
// 		log.Println("Loading plugin: ", p)
// 		o.LoadPlugin(p)
// 	}
// }

// func (o *OttO) LoadPlugin(p string) {
// 	plug, err := plugin.Open(p)
// 	if err != nil {
// 		log.Println("Failed to load plugin", err)
// 		return
// 	}

// 	ctlsym, err := plug.Lookup("Controller")
// 	if err != nil {
// 		fmt.Println("Lookup Init symbol", err)
// 		return
// 	}

// 	// 3. Assert that loaded symbol is of a desired type
// 	// in this case interface type Greeter (defined above)
// 	var c controller
// 	c, ok := ctlsym.(controller)
// 	if !ok {
// 		fmt.Println("Load Plugins: unexpected type from module symbol")
// 		return
// 	}

// 	// 4. use the module
// 	c.Init()
// }

// // Subscribe to MQTT Message the Sub (Subscriber) will is an
// // Interface that must have a Callback(topic string, payload []byte)
// // signature
// func (o *OttO) Subscribe(topic string, s Sub) {
// 	mfunc := func(c gomqtt.Client, m gomqtt.Message) {
// 		s.Callback(m.Topic(), m.Payload())
// 	}
// 	o.MQTT.Sub(topic, topic, mfunc)
// }

// func (o *OttO) Register(path string, h http.Handler) {
// 	o.Server.Register(path, h)
// }

// func (o *OttO) Publish(t string, v interface{}) {
// 	o.MQTT.Publish(t, v)
// }
