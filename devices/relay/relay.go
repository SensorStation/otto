package relay

import (
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

type Relay struct {
	*devices.DigitalDevice
}

func New(name string, offset int) *Relay {
	relay := &Relay{
		DigitalDevice: devices.NewDigitalDevice(name, offset, gpiocdev.AsOutput(0)),
	}
	return relay
}

func (r *Relay) Callback(msg *messanger.Msg) {
	str := msg.String()
	switch str {
	case "off", "0":
		r.Off()

	case "on", "1":
		r.On()
	}
	return
}
