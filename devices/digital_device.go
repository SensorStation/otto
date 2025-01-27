package devices

import (
	"fmt"

	"github.com/warthog618/go-gpiocdev"
)

// GPIODevice satisfies the Device interface and the GPIO Pin interface
type DigitalDevice struct {
	*BaseDevice
	*DigitalPin
}

// NewDigitalDevice.go creates a GPIO based device associated with GPIO Pin idx with the
// corresponding gpiocdev.LineReqOptions. If there is an gpiocdev.EventHandler an
// Event Q channel will be created.  Also a publication topic will be associated.
// For input devices the Pubs will be used to publish output state to the device.
// For output devices the Pubs will be used to publish the latest data collected
func NewDigitalDevice(name string, idx int, opts ...gpiocdev.LineReqOption) *DigitalDevice {
	d := &DigitalDevice{
		BaseDevice: NewDevice(name),
	}

	// look for an eventhandler, if so setup an event channel
	for _, opt := range opts {
		t := fmt.Sprintf("%T", opt)
		if t == "gpiocdev.EventHandler" {
			d.evtQ = make(chan gpiocdev.LineEvent)
		}
	}

	// append the pubs
	gpio = GetGPIO()
	d.DigitalPin = gpio.Pin(name, idx, opts...)
	return d
}
