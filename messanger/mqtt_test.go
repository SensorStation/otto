package messanger

import (
	"fmt"
	"testing"

	"github.com/sensorstation/otto/message"
)

type tclient struct {
	gotit bool
	topic string
	msg   string
}

func (t *tclient) Callback(msg *message.Msg) error {
	if msg.Path[0] != "t" || msg.Path[1] != "test" {
		return fmt.Errorf("bad path: %v", msg.Path)
	}

	if string(msg.Data) != "message" {
		return fmt.Errorf("bad data: %s", msg.Data)
	}
	t.gotit = true
	return nil
}

func TestSubscribe(t *testing.T) {
	c := GetMockClient()
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

	m.Subscribe(tc.topic, tc.Callback)
	m.Publish(tc.topic, tc.msg)

	if tc.gotit == false {
		t.Error("Expected to recv message but did not")
	}
}
