package otto

import (
	"testing"
)

type tclient struct {
	gotit bool
	topic string
	msg   string
}

func (t tclient) SubCallback(topic string, data []byte) {
	// Todo something
}

func TestSubscribe(t *testing.T) {
	m := GetMQTT()s
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
