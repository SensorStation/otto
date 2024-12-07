package otto

import (
	"fmt"
	"log"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

type Sub interface {
	SubCallback(topic string, data []byte)
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

func (m *MQTT) IsConnected() bool {
	if m.Client == nil {
		return false
	}
	return m.Client.IsConnected()
}

// Connect to the MQTT broker after setting some MQTT options
// then connecting to the MQTT broker
func (m *MQTT) Connect() error {

	if m.Debug {
		gomqtt.DEBUG = log.Default()
		gomqtt.ERROR = log.Default()
	}

	m.Broker = "tcp://" + m.Broker + ":1883"

	// connOpts := gomqtt.NewClientOptions().AddBroker(m.Broker).SetClientID(m.ID).SetCleanSession(true)
	opts := gomqtt.NewClientOptions()
	opts.AddBroker(m.Broker)
	opts.SetClientID(m.ID)
	opts.SetCleanSession(true)
	m.Client = gomqtt.NewClient(opts)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		l.Error("MQTT Connect: ", "error", token.Error())
		return fmt.Errorf("Failed to connect to MQTT broker %s", token.Error())
	}
	return nil
}

// Subscribe to MQTT messages that follow specific topic patterns
// wildcards '+' and '#' are supported.  Examples
// ss/<ethaddr>/<data>/tempf value
// ss/<ethaddr>/<data>/humidity value
func (m *MQTT) Sub(id string, path string, f gomqtt.MessageHandler) error {
	sub := &Subscriber{id, path, f}
	m.Subscribers[id] = sub

	if m.Client == nil {
		l.Error("MQTT Client is not connected to a broker")
		return fmt.Errorf("MQTT Client is not connected to broker: %s", m.Broker)
	}

	qos := 0
	if token := m.Client.Subscribe(path, byte(qos), f); token.Wait() && token.Error() != nil {

		// TODO: add routing that automatically subscribes subscribers when a
		// connection has been made
		return token.Error()
	} else {
		l.Debug("subscribe ", "token", token)
	}
	return nil
}

// Publish will publish a value to the given channel
func (m MQTT) Publish(topic string, value interface{}) {
	// l.Printf("[I] MQTT Publishing %s -> %v", topic, value)
	var t gomqtt.Token

	if m.Client == nil {
		l.Warn("MQTT Client is not connected to a broker")
		return
	}

	if t = m.Client.Publish(topic, byte(0), false, value); t == nil {
		if false {
			l.Info("MQTT Pub NULL token: ", "topic", topic, "value", value)
		}
		return
	}

	t.Wait()
	if t.Error() != nil {
		l.Error("MQTT Publish token: ", "error", t.Error())
	}

}

func (m *MQTT) Subscribe(topic string, s Sub) {
	mfunc := func(c gomqtt.Client, m gomqtt.Message) {
		s.SubCallback(m.Topic(), m.Payload())
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

func (mp *MQTTPrinter) SubCallback(topic string, data []byte) {
	fmt.Println(topic, " ", string(data))
}
