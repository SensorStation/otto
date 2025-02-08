package relay

import (
	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

type Relay struct {
	*device.Device
	*device.DigitalPin
}

func New(name string, offset int) *Relay {
	relay := &Relay{
		Device: device.NewDevice(name),
	}
	g := device.GetGPIO()
	relay.DigitalPin = g.Pin(name, offset, gpiocdev.AsOutput(0))
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
