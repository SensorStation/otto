package messanger

import (
	"fmt"
	"log"
	"log/slog"
	"os"

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

	gomqtt.Client `json:"-"`
}

// NewMQTT creates a new instance of the MQTT client type.
func NewMQTT() *MQTT {
	mqtt := &MQTT{
		ID:     "otto",
		Broker: "localhost",
	}
	return mqtt
}

func SetMQTTClient(c gomqtt.Client) *MQTT {
	mqtt = GetMQTT()
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

	broker := os.Getenv("MQTT_BROKER")
	if broker != "" {
		mqtt.Broker = broker
	} else {
		mqtt.Broker = "10.11.1.11"
	}
	url := "tcp://" + mqtt.Broker + ":1883"
	opts := gomqtt.NewClientOptions()
	opts.AddBroker(url)
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

// Subscribe will cause messangers to the given topic to be passed along to the
// MsgHandle f
func (m *MQTT) Subscribe(topic string, f MsgHandler) error {
	if m.Client == nil {
		slog.Error("MQTT Client is not connected to a broker")
		return fmt.Errorf("MQTT Client is not connected to broker: %s", m.Broker)
	}

	var err error
	token := m.Client.Subscribe(topic, byte(0), func(c gomqtt.Client, m gomqtt.Message) {
		slog.Debug("MQTT incoming: ", "topic", m.Topic(), "payload", string(m.Payload()))
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

// Publish will publish a value to the given topic
func (m *MQTT) Publish(topic string, value any) {
	var t gomqtt.Token

	if topic == "" {
		panic("topic is nil")
	}

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

func (m *MQTT) Close() {
	m.Client.Disconnect(1000)
}
