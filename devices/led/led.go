package led

import (
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/message"
	"github.com/warthog618/go-gpiocdev"
)

type LED struct {
	*devices.DeviceGPIO
}

func New(name string, offset int) *LED {
	led := &LED{
		DeviceGPIO: devices.NewDeviceGPIO(name, offset, devices.ModeOutput, gpiocdev.AsOutput(0)),
	}
	return led
}

func (l *LED) Callback(msg *message.Msg) error {
	switch msg.String() {
	case "off":
		l.Off()

	case "on":
		l.On()
	}
	return nil
}
