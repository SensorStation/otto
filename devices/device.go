package devices

import (
	"time"

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
	Pubs []string
	Subs []string
	Mode

	Period time.Duration
	EvtQ   chan gpiocdev.LineEvent
}

func NewDevice(name string) *Device {
	d := &Device{
		Name: name,
	}
	return d
}

func (d *Device) AddPub(p string) {
	d.Pubs = append(d.Pubs, p)
}
