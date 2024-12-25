package otto

import (
	"fmt"
	"log"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sensorstation/otto/message"
)

var (
	mqtt *MQTT
)

// MQTT is a wrapper around the Paho MQTT Go package
// Wraps the Broker, ID and Debug variables.
type MQTT struct {
	ID     string
	Broker string
	Debug  bool

	Subscribers map[string][]*Sub
	gomqtt.Client
}

// NewMQTT creates a new instance of the MQTT client type.
func NewMQTT() *MQTT {
	mqtt := &MQTT{
		ID:     "otto",
		Broker: "localhost",
	}
	mqtt.Subscribers = make(map[string][]*Sub)
	return mqtt
}

func GetMQTTClient(c gomqtt.Client) *MQTT {
	mqtt = NewMQTT()
	mqtt.Client = c
	return mqtt
}

func GetMQTT() *MQTT {
	if mqtt == nil {
		mqtt = NewMQTT()
	}
	if !mqtt.IsConnected() {
		mqtt.Connect()
	}
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
func (m *MQTT) Sub(id string, path string, f gomqtt.MessageHandler, cb Subscriber) error {
	sub := &Sub{id, path, f, cb}
	m.Subscribers[path] = append(m.Subscribers[path], sub)

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
func (mqtt *MQTT) Subscribe(topic string, s Subscriber) {
	mfunc := func(c gomqtt.Client, m gomqtt.Message) {
		msg := message.New(m.Topic(), m.Payload(), "mqtt-sub")
		for _, sub := range mqtt.Subscribers[m.Topic()] {
			sub.Subscriber.Callback(msg)
		}
	}
	mqtt.Sub(topic, topic, mfunc, s)
}

// Sub contains a Subscriber ID, a topic Path and a Message Handler
// for messages to the corresponding topic path
type Sub struct {
	ID   string
	Path string
	gomqtt.MessageHandler
	Subscriber
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

// Callback will print out all messages sent to the given topic
// from the MQTTPrinter
func (mp *MQTTPrinter) Callback(msg *message.Msg) {
	fmt.Printf("%+v\n", msg)
}
