package blasters

import (
	"fmt"
	"time"

	"github.com/sensorstation/otto"
)

// MQTTBlaster is a virtual station that will spew messages to a given
// topic to be used for testing.
type MQTTBlaster struct {
	*otto.Station
	Topic string
}

// MQTTBlasters is a collection of blaster stations that will spew
// messages from each of the different stations.
type MQTTBlasters struct {
	Count    int
	Blasters []*MQTTBlaster
	Running  bool
	Wait     int
}

// NewMQTTBlasters will create a count number of blasters ready to
// blast MQTT messages
func NewMQTTBlasters(count int) *MQTTBlasters {
	mb := &MQTTBlasters{
		Count:   count,
		Running: false,
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

// Blast will start the configured blasters to start blasting
// messages.
func (mb *MQTTBlasters) Blast() error {

	mqtt, err := otto.GetMQTT()
	if err != nil {
		return err
	}
	mqtt.Connect()
	if !mqtt.IsConnected() {
		return fmt.Errorf("MQTT Client is not connected to a broker")
	}

	wd := &WeatherData{}

	// now start blasting
	mb.Running = true
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

// Stop will cause the blasters to stop blasting.
func (mb *MQTTBlasters) Stop() {
	mb.Running = false
}
