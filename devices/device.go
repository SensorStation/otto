package devices

import (
	"time"

	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

type Mode int

const (
	ModeNone Mode = iota
	ModeInput
	ModeOutput
	ModePWM
)

type Device struct {
	Name string
	Pub  string
	Subs []string
	Mode

	Period time.Duration
	EvtQ   chan gpiocdev.LineEvent
}

func NewDevice(name string) *Device {
	d := &Device{
		Name: name,
	}
	// should we do this here?
	d.AddPub(name)
	return d
}

func (d *Device) AddPub(p string) {
	d.Pub = p
}

func (d *Device) Publish(data any) {
	m := messanger.GetMQTT()
	if d.Pub == "" {
		logger.GetLogger().Error("Device.Publish failed has no pub", "name", d.Name)
	}
	m.Publish(d.Pub, data)
}
