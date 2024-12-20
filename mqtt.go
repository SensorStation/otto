package otto

import (
	"fmt"
	"log"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

// Subscriber is an interface that defines a struct needs to have the
// SubCallback(topic string, data []byte) function defined.
type Subscriber interface {
	// SubCallback(topic string, data []byte)
	SubCallback(msg *Msg)
}

// MQTT is a wrapper around the Paho MQTT Go package
// Wraps the Broker, ID and Debug variables.
type MQTT struct {
	ID     string
	Broker string
	Debug  bool

	Subscribers map[string]*Sub
	gomqtt.Client
}

// NewMQTT creates a new instance of the MQTT client type.
func NewMQTT() *MQTT {
	mqtt := &MQTT{
		ID:     "otto",
		Broker: "localhost",
	}
	mqtt.Subscribers = make(map[string]*Sub)
	return mqtt
}

// IsConnected will tell you if the MQTT client is connected to
// the specified broker
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

	// If we are testing m.Client will point to the mock client otherwise
	// in real life a new real client will be created
	if m.Client == nil {
		m.Client = gomqtt.NewClient(opts)
	}

	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		l.Error("MQTT Connect: ", "error", token.Error())
		return fmt.Errorf("Failed to connect to MQTT broker %s", token.Error())
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

// Subscribe to MQTT messages that follow specific topic patterns
// wildcards '+' and '#' are supported.  Examples
// ss/<ethaddr>/<data>/tempf value
// ss/<ethaddr>/<data>/humidity value
func (m *MQTT) Sub(id string, path string, f gomqtt.MessageHandler) error {
	sub := &Sub{id, path, f}
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

// Subscribe causes the MQTT client to subscribe to the given topic with
// the connected broker
func (m *MQTT) Subscribe(topic string, s Subscriber) {
	mfunc := func(c gomqtt.Client, m gomqtt.Message) {
		// s.SubCallback(m.Topic(), m.Payload())
		msg := NewMsg(m.Topic(), m.Payload(), "mqtt-sub")
		s.SubCallback(msg)
	}
	m.Sub(topic, topic, mfunc)
}

// Sub contains a Subscriber ID, a topic Path and a Message Handler
// for messages to the corresponding topic path
type Sub struct {
	ID   string
	Path string
	gomqtt.MessageHandler
}

// String returns a string representation of the Subscriber and
// Subscriber ID
func (sub *Sub) String() string {
	return sub.ID + " " + sub.Path
}

// MQTTPrinter defines the struct that simply prints what ever
// message is sent to a given topic
type MQTTPrinter struct {
}

// SubCallback will print out all messages sent to the given topic
// from the MQTTPrinter
func (mp *MQTTPrinter) SubCallback(msg *Msg) {
	fmt.Printf("%+v\n", msg)
}
