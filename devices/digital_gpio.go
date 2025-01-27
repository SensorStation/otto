package devices

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/warthog618/go-gpiocdev"
)

// GPIO is used to initialize the GPIO and pins on a raspberry pi
type GPIO struct {
	Chipname string              `json:"chipname"`
	pins     map[int]*DigitalPin `json:"pins"`
	Mock     bool                `json:"mock"`
}

var (
	gpio *GPIO
)

// GetGPIO returns the GPIO singleton for the Raspberry PI
func GetGPIO() *GPIO {
	if gpio == nil {
		gpio = &GPIO{
			Chipname: "gpiochip4", // raspberry pi-5
		}
		gpio.pins = make(map[int]*DigitalPin)
	}
	return gpio
}

// Init initialized the GPIO
func (gpio *GPIO) Init() error {
	for _, pin := range gpio.pins {
		if err := pin.Init(); err != nil {
			slog.Error("Error initializing pin ", "offset", pin.offset)
		}
	}
	return nil
}

// Pin initializes the given GPIO pin, name and mode
func (gpio *GPIO) Pin(name string, offset int, opts ...gpiocdev.LineReqOption) (p *DigitalPin) {
	p = &DigitalPin{
		offset: offset,
		opts:   opts,
	}

	if gpio.pins == nil {
		gpio.pins = make(map[int]*DigitalPin)
	}
	gpio.pins[offset] = p
	if err := p.Init(); err != nil {
		slog.Error(err.Error(), "name", name, "offset", offset)
	}
	return p
}

// Shutdown resets the GPIO line allowing use by another program
func (gpio *GPIO) Close() {
	for _, p := range gpio.pins {
		p.Reconfigure(gpiocdev.AsInput)
		p.Close()
	}
	gpio.pins = nil
}

// String returns the string representation of the GPIO
func (gpio *GPIO) String() string {
	str := ""
	for _, pin := range gpio.pins {
		str += pin.String()
	}
	return str
}

// JSON returns the JSON representation of the GPIO
func (gpio *GPIO) JSON() (j []byte, err error) {
	j, err = json.Marshal(gpio)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling GPIO %s", err)
	}
	return j, nil
}
