package otto

import (
	"testing"

	"github.com/sensorstation/otto/message"
	"github.com/sensorstation/otto/mocks"
)

type tclient struct {
	gotit bool
	topic string
	msg   string
}

func (t *tclient) SubCallback(msg *message.Msg) {
	println(msg.Path)
	if msg.Path[0] != "t" || msg.Path[1] != "test" {
		return
	}

	println(msg.Data)
	if string(msg.Data) != "message" {
		return
	}
	t.gotit = true
}

func TestSubscribe(t *testing.T) {
	c := mocks.GetMockClient()
	m := GetMQTTClient(c)
	err := m.Connect()
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
