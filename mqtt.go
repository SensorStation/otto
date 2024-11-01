package otto

import (
	"fmt"
	"log"
	"os"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTT is a wrapper around the Paho MQTT Go package
// Wraps the Broker, ID and Debug variables.
type MQTT struct {
	ID     string
	Broker string
	Debug  bool

	subscribers map[string]*Subscriber
	gomqtt.Client
}

// Start causes a Connect to the MQTT Broker
func (m *MQTT) Start() {
	m.subscribers = make(map[string]*Subscriber)
	m.Connect()
}

// Connect to the MQTT broker after setting some MQTT options
// then connecting to the MQTT broker
func (m *MQTT) Connect() {
	if m.Debug {
		gomqtt.DEBUG = log.New(os.Stdout, "", 0)
		gomqtt.ERROR = log.New(os.Stdout, "", 0)
	}

	m.Broker = "tcp://" + m.Broker + ":1883"

	// connOpts := gomqtt.NewClientOptions().AddBroker(m.Broker).SetClientID(m.ID).SetCleanSession(true)
	opts := gomqtt.NewClientOptions()
	opts.AddBroker(m.Broker)
	opts.SetClientID(m.ID)
	opts.SetCleanSession(true)
	m.Client = gomqtt.NewClient(opts)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("MQTT Connect: ", token.Error())
		return
	}
	log.Println("Connected to broker: ", m.Broker)
}

// Subscribe to MQTT messages that follow specific topic patterns
// wildcards '+' and '#' are supported.  Examples
//
// ss/<ethaddr>/<data>/tempf value
// ss/<ethaddr>/<data>/humidity value
func (m *MQTT) Subscribe(id string, path string, f gomqtt.MessageHandler) {
	sub := &Subscriber{id, path, f}
	m.subscribers[id] = sub

	qos := 0
	if token := m.Client.Subscribe(path, byte(qos), f); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		if m.Debug {
			log.Printf("subscribe token: %v", token)
		}
	}
	log.Println(id, "subscribed to", path)
}

// Subscriber contains a Subscriber ID, a topic Path and a
// Message Handler for messages to the corresponding topic path
type Subscriber struct {
	ID   string
	Path string
	gomqtt.MessageHandler
}

// String returns a string representation of the Subscriber and
// Subscriber ID
func (sub *Subscriber) String() string {
	return sub.ID + " " + sub.Path
}

// Publish will publish a value to the given channel
func (m MQTT) Publish(topic string, value interface{}) {
	log.Printf("[I] MQTT Publishing %s -> %v", topic, value)
	var t gomqtt.Token

	if t = m.Client.Publish(topic, byte(0), false, value); t == nil {
		if true {
			log.Printf("[I] MQTT Pub NULL token: %s - %v", topic, value)
		}
		return
	}
	t.Wait()
	if t.Error() != nil {
		fmt.Printf("MQTT Publish token: %+v\n", t.Error())
	}

}
