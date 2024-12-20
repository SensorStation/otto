package blasters

import (
	"fmt"
	"time"

	"github.com/sensorstation/otto"
)

type MQTTBlaster struct {
	*otto.Station
	Topic string
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
		Wait:    2000,
	}

	mb.Blasters = make([]*MQTTBlaster, mb.Count)
	for i := 0; i < mb.Count; i++ {
		id := fmt.Sprintf("station-%d", i)
		topic := fmt.Sprintf("ss/d/%s/temphum", id)

		mb.Blasters[i] = &MQTTBlaster{
			Topic:   topic,
			Station: otto.NewStation(id),
		}
	}
	return mb
}

func (mb *MQTTBlasters) Blast() error {

	mqtt, err := otto.GetMQTT()
	if err != nil {
		return err
	}
	if !mqtt.IsConnected() {
		return fmt.Errorf("MQTT Client is not connected to a broker")
	}

	wd := &WeatherData{}

	// now start blasting
	for mb.Running {
		for i := 0; i < mb.Count; i++ {
			b := mb.Blasters[i]

			msg := wd.NewMsg()
			mqtt.Publish(b.Topic, msg.Byte())
		}
		time.Sleep(time.Duration(mb.Wait) * time.Millisecond)
	}
	otto.GetLogger().Info("MQTT Blaster has stopped")
	return nil
}

func (mb *MQTTBlasters) Stop() {
	mb.Running = false
}
