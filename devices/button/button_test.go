package button

import (
	"strconv"
	"testing"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/message"
	"github.com/sensorstation/otto/mocks"
)

var (
	gotit [2]bool
)

func TestButton(t *testing.T) {
	done := make(chan bool)

	c := mocks.GetMockClient()
	m := otto.GetMQTTClient(c)
	err := m.Connect()
	if err != nil {
		t.Error("Failed to connect to MQTT broker: ", err)
	}

	devices.GetGPIO().Mock = true

	b := New("button", 23)
	m.Subscribe("ss/c/"+otto.StationName+"/button", b)
	go b.EventLoop(done)
	b.Line.(*devices.MockLine).MockHWInput(0)
	b.Line.(*devices.MockLine).MockHWInput(1)
	done <- true

	if !gotit[0] || !gotit[1] {
		t.Errorf("failed to get 0 and 1 got (%t) and (%t)", gotit[0], gotit[1])
	}

}

func (b *Button) Callback(msg *message.Msg) {
	i, err := strconv.Atoi(string(msg.Data))
	if err != nil {
		return
	}
	gotit[i] = true
}
