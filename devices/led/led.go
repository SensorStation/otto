package led

import (
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/message"
)

type LED struct {
	*devices.GPIODevice
}

func New(name string, offset int) *LED {
	led := &LED{
		GPIODevice: devices.GPIOOut(name, offset),
	}
	return led
}

func (l *LED) Callback(msg *message.Msg) {
	switch msg.Path[3] {
	case "off":
		l.GPIODevice.Off()

	case "on":
		l.On()
	}
}
