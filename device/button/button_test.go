package button

import (
	"strconv"
	"testing"
	"time"

	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/messanger"
)

var (
	gotit [2]bool
)

func TestButton(t *testing.T) {
	device.Mock(true)
	done := make(chan any)

	c := messanger.GetMockClient()
	m := messanger.SetMQTTClient(c)
	err := m.Connect()
	if err != nil {
		t.Error("Failed to connect to MQTT broker: ", err)
	}

	b := New("button", 23)
	b.AddPub(messanger.TopicControl("button"))
	b.Subscribe(messanger.TopicControl("button"), b.Callback)
	go b.EventLoop(done, b.ReadPub)

	b.MockHWInput(0)
	b.MockHWInput(1)

	// This is weak figure out a better solution
	time.Sleep(10 * time.Millisecond)
	b.Close()
	done <- true

	if !gotit[0] || !gotit[1] {
		t.Errorf("failed to get 0 and 1 got (%t) and (%t)", gotit[0], gotit[1])
	}

}

func (b *Button) Callback(msg *messanger.Msg) {
	i, err := strconv.Atoi(msg.String())
	if err != nil {
		return
	}
	gotit[i] = true
	return
}
