package msg

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Subcriber represent MQTT subscriptions allow us to recieve data
// from MQTT. This data will end up being passed along to the
// Consumer interface

type Subscriber struct {
	ID string
	Path string
	mqtt.MessageHandler
	Consumers []Consumer
}

func (sub *Subscriber) String() string {
	return sub.ID + " " + sub.Path
}

