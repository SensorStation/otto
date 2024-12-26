package relay

import (
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/message"
)

type Relay struct {
	*devices.DeviceGPIO
}

func New(name string, offset int) *Relay {
	relay := &Relay{
		DeviceGPIO: devices.NewDeviceGPIO(name, offset, devices.ModeOutput),
	}
	return relay
}

func (r *Relay) Callback(msg *message.Msg) {
	switch msg.Path[3] {
	case "off":
		r.Off()

	case "on":
		r.On()
	}
}
