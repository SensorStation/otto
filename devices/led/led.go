package led

import (
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

type LED struct {
	*devices.DigitalDevice
}

func New(name string, offset int) *LED {
	led := &LED{
		DigitalDevice: devices.NewDigitalDevice(name, offset, gpiocdev.AsOutput(0)),
	}
	return led
}

func (l *LED) Callback(msg *messanger.Msg) {
	switch msg.String() {
	case "off", "0":
		l.Off()

	case "on", "1":
		l.On()
	}
	return
}
