package led

import (
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/message"
	"github.com/warthog618/go-gpiocdev"
)

type LED struct {
	*devices.GPIODevice
}

func New(name string, offset int) *LED {
	led := &LED{
		GPIODevice: devices.NewGPIODevice(name, offset, devices.ModeOutput, gpiocdev.AsOutput(0)),
	}
	return led
}

func (l *LED) Callback(msg *message.Msg) {
	switch msg.String() {
	case "off":
		l.Off()

	case "on":
		l.On()
	}
	return
}
