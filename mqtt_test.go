package otto

import (
	"testing"
	"time"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

type MockClient struct {
	connected bool
}

func (m MockClient) IsConnected() bool {
	return m.connected
}

func (m MockClient) IsConnectionOpen() bool {
	return m.connected
}

func (m MockClient) Connect() gomqtt.Token {
	var t MockToken
	return t
}

func (m MockClient) Disconnect(quiecense uint) {
}

func (m MockClient) Publish(topic string, qos byte, retained bool, payload interface{}) gomqtt.Token {
	var t MockToken
	return t
}

func (m MockClient) Subscribe(topic string, qos byte, callback gomqtt.MessageHandler) gomqtt.Token {
	var t MockToken
	return t
}

func (m MockClient) SubscribeMultiple(filters map[string]byte, callback gomqtt.MessageHandler) gomqtt.Token {
	var t MockToken
	return t
}

func (m MockClient) Unsubscribe(topics ...string) gomqtt.Token {
	var t MockToken
	return t
}

func (m MockClient) AddRoute(topic string, callback gomqtt.MessageHandler) {

}

func (m MockClient) OptionsReader() gomqtt.ClientOptionsReader {
	var r gomqtt.ClientOptionsReader
	return r
}

type MockToken struct {
}

func (t MockToken) Wait() bool {
	return true
}

func (t MockToken) WaitTimeout(d time.Duration) bool {
	return true
}

func (t MockToken) Done() <-chan struct{} {
	return make(chan struct{})
}

func (t MockToken) Error() error {
	return nil
}

type tclient struct {
	gotit bool
	topic string
	msg   string
}

func (t tclient) SubCallback(topic string, data []byte) {
	// Todo something
}

func TestSubscribe(t *testing.T) {
	m := GetMQTT()
	m.Client = MockClient{}

	err := m.Connect()
	if err != nil {
		t.Error("Failed to connect to MQTT broker: ", err)
	}

	tc := &tclient{
		gotit: true,
		topic: "t/test",
		msg:   "message",
	}

	m.Subscribe(tc.topic, tc)
	m.Publish(tc.topic, tc.msg)

	if tc.gotit == false {
		t.Error("Expected to recv message but did not")
	}
}
