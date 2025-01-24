package devices

import (
	"fmt"

	"github.com/warthog618/go-gpiocdev"
)

// GPIODevice satisfies the Device interface and the GPIO Pin interface
type GPIODevice struct {
	*Device
	*Pin
}

// NewGPIODevice creates a GPIO based device associated with GPIO Pin idx with the
// corresponding gpiocdev.LineReqOptions. If there is an gpiocdev.EventHandler an
// Event Q channel will be created.  Also a publication topic will be associated.
// For input devices the Pubs will be used to publish output state to the device.
// For output devices the Pubs will be used to publish the latest data collected
func NewGPIODevice(name string, idx int, mode Mode, opts ...gpiocdev.LineReqOption) *GPIODevice {
	d := &GPIODevice{
		Device: NewDevice(name),
	}

	// look for an eventhandler, if so setup an event channel
	for _, opt := range opts {
		t := fmt.Sprintf("%T", opt)
		if t == "gpiocdev.EventHandler" {
			d.EvtQ = make(chan gpiocdev.LineEvent)
		}
	}

	// append the pubs
	gpio = GetGPIO()
	d.Pin = gpio.Pin(name, idx, opts...)
	return d
}
