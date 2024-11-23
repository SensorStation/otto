package otto

import (
	"fmt"
	"log"
	"os"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

type Sub interface {
	Callback(msg *Msg)
}

// MQTT is a wrapper around the Paho MQTT Go package
// Wraps the Broker, ID and Debug variables.
type MQTT struct {
	ID     string
	Broker string
	Debug  bool

	Subscribers map[string]*Subscriber
	gomqtt.Client
}

func NewMQTT() *MQTT {
	mqtt := &MQTT{
		ID:     "otto",
		Broker: "localhost",
	}
	mqtt.Subscribers = make(map[string]*Subscriber)
	return mqtt
}

// Connect to the MQTT broker after setting some MQTT options
// then connecting to the MQTT broker
func (m *MQTT) Connect() error {
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
		return fmt.Errorf("Failed to connect to MQTT broker %s", token.Error())
	}
	return nil
}

// Subscribe to MQTT messages that follow specific topic patterns
// wildcards '+' and '#' are supported.  Examples
// ss/<ethaddr>/<data>/tempf value
// ss/<ethaddr>/<data>/humidity value
func (m *MQTT) Sub(id string, path string, f gomqtt.MessageHandler) {
	sub := &Subscriber{id, path, f}
	m.Subscribers[id] = sub

	qos := 0
	if token := m.Client.Subscribe(path, byte(qos), f); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		if m.Debug {
			log.Printf("subscribe token: %v", token)
		}
	}
}

// Publish will publish a value to the given channel
func (m MQTT) Publish(topic string, value interface{}) {
	// log.Printf("[I] MQTT Publishing %s -> %v", topic, value)
	var t gomqtt.Token

	if t = m.Client.Publish(topic, byte(0), false, value); t == nil {
		if false {
			log.Printf("[I] MQTT Pub NULL token: %s - %v", topic, value)
		}
		return
	}

	t.Wait()
	if t.Error() != nil {
		fmt.Printf("MQTT Publish token: %+v\n", t.Error())
	}

}

func (m *MQTT) Subscribe(topic string, s Sub) {
	mfunc := func(c gomqtt.Client, m gomqtt.Message) {

		// MQTT Middleware here
		msg, err := MsgFromMQTT(m.Topic(), m.Payload())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse mqtt message topic: %s message: %s - err %s\n",
				topic, string(m.Payload()), err)
			return
		}
		msg.Source = "mqtt"
		s.Callback(msg)
	}
	m.Sub(topic, topic, mfunc)
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

type MQTTPrinter struct {
}

func (mp *MQTTPrinter) Callback(msg *Msg) {
	fmt.Printf("  ID: %d\n", msg.ID)
	fmt.Printf("Path: %q\n", msg.Path)
	fmt.Printf("Args: %q\n", msg.Args)
	fmt.Printf(" Msg: %s\n", string(msg.Message))
	fmt.Printf(" Src: %s\n", msg.Source)
	fmt.Printf("Time: %s\n", msg.Time)
	fmt.Println()
}
