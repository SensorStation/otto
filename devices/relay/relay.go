package relay

import (
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/message"
	"github.com/warthog618/go-gpiocdev"
)

type Relay struct {
	*devices.GPIODevice
}

func New(name string, offset int) *Relay {
	relay := &Relay{
		GPIODevice: devices.NewGPIODevice(name, offset, devices.ModeOutput, gpiocdev.AsOutput(0)),
	}
	return relay
}

func (r *Relay) Callback(msg *message.Msg) {
	str := msg.String()
	switch str {
	case "off":
		r.Off()

	case "on":
		r.On()
	}
	return
}
