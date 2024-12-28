package devices

import (
	"fmt"

	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

type DeviceGPIO struct {
	*Device
	*Pin
}

func NewDeviceGPIO(name string, idx int, mode Mode, opts ...gpiocdev.LineReqOption) *DeviceGPIO {
	d := &DeviceGPIO{
		Device: &Device{
			Name: name,
			Mode: mode,
		},
	}

	for _, opt := range opts {
		t := fmt.Sprintf("%T", opt)
		if t == "gpiocdev.EventHandler" {
			d.EvtQ = make(chan gpiocdev.LineEvent)
		}
	}

	d.Pubs = append(d.Pubs, messanger.TopicData(name))
	gpio = GetGPIO()
	d.Pin = gpio.Pin(name, idx, opts...)
	return d
}

func (d *DeviceGPIO) On() {
	d.Set(1)
}

func (d *DeviceGPIO) Off() {
	d.Set(0)
}

func (d *DeviceGPIO) Set(v int) {
	d.Pin.Set(v)
	val := "off"
	if v > 0 {
		val = "on"
	}

	m := messanger.GetMQTT()
	for _, p := range d.Pubs {
		m.Publish(p, val)
	}
}
