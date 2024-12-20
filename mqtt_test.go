package otto

import (
	"testing"
	"time"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

type MockMessage struct {
	topic   string
	payload []byte
}

func (m *MockMessage) Duplicate() bool {
	return false
}

func (m *MockMessage) Qos() byte {
	return 0
}

func (m *MockMessage) Retained() bool {
	return false
}

func (m *MockMessage) Topic() string {
	return m.topic
}

func (m *MockMessage) MessageID() uint16 {
	return 10
}

func (m *MockMessage) Payload() []byte {
	return m.payload
}

func (m *MockMessage) Ack() {
}

type MockClient struct {
	connected bool
	mqtt      *MQTT
}

func GetMockMQTT() (*MQTT, error) {
	mqtt := NewMQTT()
	mqtt.Client = MockClient{
		mqtt: mqtt,
	}
	return mqtt, nil
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
	if sub, ex := m.mqtt.Subscribers[topic]; ex {
		mes := &MockMessage{
			topic:   topic,
			payload: []byte(payload.(string)),
		}
		sub.MessageHandler(m, mes)
	}
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

func (t *tclient) SubCallback(msg *Msg) {
	println(msg.Path)
	if msg.Path[0] != "t" || msg.Path[1] != "test" {
		return
	}

	println(msg.Message)
	if string(msg.Message) != "message" {
		return
	}
	t.gotit = true
}

func TestSubscribe(t *testing.T) {
	m, err := GetMockMQTT()
	if err != nil {
		t.Logf(err.Error())
		return
	}

	err = m.Connect()
	if err != nil {
		t.Error("Failed to connect to MQTT broker: ", err)
	}

	tc := &tclient{
		gotit: false,
		topic: "t/test",
		msg:   "message",
	}

	m.Subscribe(tc.topic, tc)
	m.Publish(tc.topic, tc.msg)

	if tc.gotit == false {
		t.Error("Expected to recv message but did not")
	}
}
