package led

import (
	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/device/drivers"
	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

type LED struct {
	*device.Device
	*drivers.DigitalPin
}

func New(name string, offset int) *LED {
	led := &LED{
		Device: device.NewDevice(name),
	}
	g := drivers.GetGPIO()
	led.DigitalPin = g.Pin(name, offset, gpiocdev.AsOutput(0))
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
