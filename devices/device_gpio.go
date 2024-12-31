package devices

import (
	"fmt"

	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

// DeviceGPIO satisfies the Device interface and the GPIO Pin interface
type DeviceGPIO struct {
	*Device
	*Pin
}

// NewDeviceGPIO creates a GPIO based device associated with GPIO Pin idx with the
// corresponding gpiocdev.LineReqOptions. If there is an gpiocdev.EventHandler an
// Event Q channel will be created.  Also a publication topic will be associated.
// For input devices the Pubs will be used to publish output state to the device.
// For output devices the Pubs will be used to publish the latest data collected
func NewDeviceGPIO(name string, idx int, mode Mode, opts ...gpiocdev.LineReqOption) *DeviceGPIO {
	d := &DeviceGPIO{
		Device: NewDevice(name),
	}
	// Do we need mode?
	d.Mode = mode

	// look for an eventhandler, if so setup an event channel
	for _, opt := range opts {
		t := fmt.Sprintf("%T", opt)
		if t == "gpiocdev.EventHandler" {
			d.EvtQ = make(chan gpiocdev.LineEvent)
		}
	}

	// append the pubs
	d.AddPub(name)
	gpio = GetGPIO()
	d.Pin = gpio.Pin(name, idx, opts...)
	return d
}

// On sets the output state of the Pin to ON (1)
func (d *DeviceGPIO) On() {
	d.Set(1)
}

// Off sets the output state of the Pin to OFF (0)
func (d *DeviceGPIO) Off() {
	d.Set(0)
}

// Set the output state of the pin to the value of v Off(0) or On != 0
func (d *DeviceGPIO) Set(v int) {
	d.Pin.Set(v)
	val := "off"
	if v > 0 {
		val = "on"
	}

	m := messanger.GetMQTT()
	m.Publish(d.Pub, val)
}
