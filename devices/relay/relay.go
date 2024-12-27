package relay

import (
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/message"
	"github.com/warthog618/go-gpiocdev"
)

type Relay struct {
	*devices.DeviceGPIO
}

func New(name string, offset int) *Relay {
	relay := &Relay{
		DeviceGPIO: devices.NewDeviceGPIO(name, offset, devices.ModeOutput, gpiocdev.AsOutput(0)),
	}
	return relay
}

func (r *Relay) Callback(msg *message.Msg) {
	switch msg.String() {
	case "off":
		r.Off()

	case "on":
		r.On()
	}
}
