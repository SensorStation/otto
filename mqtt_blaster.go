package otto

import (
	"fmt"
	"time"
)

type MQTTMsg struct {
	Topic   string
	Message string
}

type MQTTBlaster struct {
	*Station
	MQTTMsg
}

type MQTTBlasters struct {
	Count    int
	Blasters []*MQTTBlaster
	Running  bool
	Wait     int
}

func NewMQTTBlasters(count int) *MQTTBlasters {
	mb := &MQTTBlasters{
		Count:   count,
		Running: true,
		Wait:    500,
	}

	mb.Blasters = make([]*MQTTBlaster, mb.Count)
	for i := 0; i < mb.Count; i++ {
		topic := fmt.Sprintf("ss/d/%d/temphum", i)
		msg := `{ "tempc": 100, "humidty": 78 }`

		id := fmt.Sprintf("station-%d", i)
		mb.Blasters[i] = &MQTTBlaster{
			Station: NewStation(id),
			MQTTMsg: MQTTMsg{
				Topic:   topic,
				Message: msg,
			},
		}
	}
	return mb
}

func (mb *MQTTBlasters) Blast() error {

	mqtt := GetMQTT()

	// now start blasting
	for mb.Running {
		for i := 0; i < mb.Count; i++ {
			b := mb.Blasters[i]
			mqtt.Publish(b.Topic, b.Message)
		}
		time.Sleep(time.Duration(mb.Wait) * time.Millisecond)
	}
	return nil
}

func (mb *MQTTBlasters) Stop() {
	mb.Running = false
}
