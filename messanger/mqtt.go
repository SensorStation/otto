package messanger

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqtt *MQTT
)

// MQTT is a wrapper around the Paho MQTT Go package
// Wraps the Broker, ID and Debug variables.
type MQTT struct {
	ID     string `json:"id"`
	Broker string `json:"broker"`
	Debug  bool   `json:"debug"`

	Subscribers map[string][]MsgHandle `json:"subscribers"`
	Publishers  map[string]int         `json:"publishers"`

	gomqtt.Client `json:"-"`
}

// NewMQTT creates a new instance of the MQTT client type.
func NewMQTT() *MQTT {
	mqtt := &MQTT{
		ID:     "otto",
		Broker: "localhost",
	}
	mqtt.Subscribers = make(map[string][]MsgHandle)
	mqtt.Publishers = make(map[string]int)

	return mqtt
}

func SetMQTTClient(c gomqtt.Client) *MQTT {
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
		slog.Error("MQTT Connect: ", "error", token.Error())
		return fmt.Errorf("Failed to connect to MQTT broker %s", token.Error())
	}
	return nil
}

// Publish messanger will publish a messanger
func (m MQTT) PublishMsg(msg *Msg) {
	m.Publish(msg.Topic, msg.Data)
}

// Publish will publish a value to the given channel
func (m MQTT) Publish(topic string, value interface{}) {
	var t gomqtt.Token

	if topic == "" {
		panic("topic is nil")
	}

	m.Publishers[topic] += 1
	if m.Client == nil {
		slog.Warn("MQTT Client is not connected to a broker")
		return
	}

	if t = m.Client.Publish(topic, byte(0), false, value); t == nil {
		if false {
			slog.Info("MQTT Pub NULL token: ", "topic", topic, "value", value)
		}
		return
	}

	t.Wait()
	if t.Error() != nil {
		slog.Error("MQTT Publish token: ", "error", t.Error())
	}

}

// Subscribe will cause messangers to the given topic to be passed along to the
// MsgHandle f
func (m *MQTT) Subscribe(topic string, f MsgHandle) error {
	m.Subscribers[topic] = append(m.Subscribers[topic], f)
	if m.Client == nil {
		slog.Error("MQTT Client is not connected to a broker")
		return fmt.Errorf("MQTT Client is not connected to broker: %s", m.Broker)
	}

	var err error
	token := m.Client.Subscribe(topic, byte(0), func(c gomqtt.Client, m gomqtt.Message) {
		slog.Info("MQTT incoming: ", "topic", m.Topic(), "payload", string(m.Payload()))
		msg := New(m.Topic(), m.Payload(), "mqtt-sub")
		f(msg)
	})

	if token.Wait() && token.Error() != nil {
		// TODO: add routing that automatically subscribes subscribers when a
		// connection has been made
		return token.Error()
	}
	return err
}

func (mqtt MQTT) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mq := struct {
		ID     string `json:"id"`
		Broker string `json:"broker"`
		Debug  bool   `json:"debug"`

		Subscribers []string
		Publishers  []string
	}{
		ID:     mqtt.ID,
		Broker: mqtt.Broker,
		Debug:  mqtt.Debug,
	}
	for s, _ := range mqtt.Subscribers {
		mq.Subscribers = append(mq.Subscribers, s)
	}
	for p, _ := range mqtt.Publishers {
		mq.Publishers = append(mq.Publishers, p)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(mq)
	if err != nil {
		slog.Error("MQTT.ServeHTTP failed to encode", "error", err)
	}
}

// MQTTPrinter defines the struct that simply prints what ever
// message is sent to a given topic
type MQTTPrinter struct {
}

// Callback will print out all messages sent to the given topic
// from the MQTTPrinter
func (mp *MQTTPrinter) Callback(msg *Msg) {
	fmt.Printf("%+v\n", msg)
	return
}
