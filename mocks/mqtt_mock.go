package mocks

import (
	"errors"
	"time"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sensorstation/otto/message"
)

type Callback func(msg *message.Msg)

type Sub struct {
	ID    string
	Topic string
	Callback
}

type MockMessage struct {
	topic   string
	payload []byte
}

func (m MockMessage) Duplicate() bool {
	return false
}

func (m MockMessage) Qos() byte {
	return 0
}

func (m MockMessage) Retained() bool {
	return false
}

func (m MockMessage) Topic() string {
	return m.topic
}

func (m MockMessage) MessageID() uint16 {
	return 10
}

func (m MockMessage) Payload() []byte {
	return m.payload
}

func (m MockMessage) Ack() {
}

type MockClient struct {
	connected   bool
	Subscribers map[string]Sub
}

func GetMockClient() *MockClient {
	cli := &MockClient{}
	return cli
}

func (m MockClient) IsConnected() bool {
	return m.connected
}

func (m MockClient) IsConnectionOpen() bool {
	return m.connected
}

func (m MockClient) Connect() gomqtt.Token {
	m.connected = true
	var t MockToken
	return t
}

func (m MockClient) Disconnect(quiecense uint) {
}

func (m MockClient) MessageHandler(c gomqtt.Client, mm gomqtt.Message) {
	// t := mm.Topic()
	// p := mm.Payload()

	// fmt.Printf("%s -> %+v\n", t, p)
}

func (m MockClient) Publish(topic string, qos byte, retained bool, payload interface{}) gomqtt.Token {
	n := root.lookup(topic)
	var t MockToken
	if n == nil {
		t.Err = errors.New("no subscribers for message")
		return t
	}
	mm := MockMessage{
		topic:   topic,
		payload: []byte(payload.(string)),
	}
	n.pub(m, mm)
	return t
}

func (m MockClient) Subscribe(topic string, qos byte, mh gomqtt.MessageHandler) gomqtt.Token {
	root.insert(topic, mh)

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
	Err error
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
