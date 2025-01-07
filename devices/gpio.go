package devices

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

// Line interface is used to emulate a GPIO pin as
// implemented by the go-gpiocdev package
type Line interface {
	Close() error
	Offset() int
	SetValue(int) error
	Reconfigure(...gpiocdev.LineConfigOption) error
	Value() (int, error)
}

// Pin represents a single GPIO Pin
type Pin struct {
	Opts []gpiocdev.LineReqOption
	Line

	offset int
	val    int
	mock   bool
}

var (
	gpio *GPIO
)

// Init the pin from the offset and mode
func (p *Pin) Init() error {

	gpio := GetGPIO()
	if gpio.Mock {
		line := GetMockLine(p.offset, p.Opts...)
		p.mock = true
		p.Line = line
		return nil
	}

	line, err := gpiocdev.RequestLine(gpio.Chipname, p.offset, p.Opts...)
	if err != nil {
		return err
	}
	p.Line = line
	return nil
}

// String returns a string representation of the GPIO pin
func (pin *Pin) String() string {
	v, err := pin.Value()
	if err != nil {
		slog.Error("Failed getting the value of ", "pin", pin.offset, "error", err)
	}
	str := fmt.Sprintf("%d: %d\n", pin.offset, v)
	return str
}

// Get returns the value of the pin, an error is returned if
// the GPIO value fails. Note: you can Get() the value of an
// input pin so no direction checks are done
func (pin *Pin) Get() (int, error) {
	if pin.Line == nil {
		return 0, fmt.Errorf("GPIO not active")
	}
	return pin.Line.Value()
}

// Set the value of the pin. Note: you can NOT set the value
// of an input pin, so we will check it and return an error.
// This maybe worthy of making it a panic!
func (pin *Pin) Set(v int) error {
	if pin.Line == nil {
		return fmt.Errorf("GPIO not active")
	}
	pin.val = v
	return pin.Line.SetValue(v)
}

// On sets the value of the pin to 1
func (pin *Pin) On() error {
	return pin.Set(1)
}

// Off sets the value of the pin to 0
func (pin *Pin) Off() error {
	return pin.Set(0)
}

// Toggle with flip the value of the pin from 1 to 0 or 0 to 1
func (pin *Pin) Toggle() error {
	val, err := pin.Get()
	if err != nil {
		return err
	}

	if val == 0 {
		val = 1
	} else {
		val = 0
	}
	return pin.Set(val)
}

// Callback is the default callback for pins if they are
// registered with the MQTT.Subscribe() function
func (pin Pin) Callback(msg *messanger.Msg) {
	cmd := msg.String()
	switch cmd {
	case "on":
		pin.On()

	case "off":
		pin.Off()

	case "toggle":
		pin.Toggle()

	}
	return
}

// GPIO is used to initialize the GPIO and pins on a raspberry pi
type GPIO struct {
	Chipname string       `json:"chipname"`
	Pins     map[int]*Pin `json:"pins"`
	Mock     bool         `json:"mock"`
}

// GetGPIO returns the GPIO singleton for the Raspberry PI
func GetGPIO() *GPIO {
	if gpio == nil {
		gpio = &GPIO{
			Chipname: "gpiochip4", // raspberry pi-5
		}
		gpio.Pins = make(map[int]*Pin)
	}
	return gpio
}

// Init initialized the GPIO
func (gpio *GPIO) Init() error {
	for _, pin := range gpio.Pins {
		if err := pin.Init(); err != nil {
			slog.Error("Error initializing pin ", "offset", pin.offset)
		}
	}
	return nil
}

// Pin initializes the given GPIO pin, name and mode
func (gpio *GPIO) Pin(name string, offset int, opts ...gpiocdev.LineReqOption) (p *Pin) {
	p = &Pin{
		offset: offset,
		Opts:   opts,
	}

	if gpio.Pins == nil {
		gpio.Pins = make(map[int]*Pin)
	}
	gpio.Pins[offset] = p
	if err := p.Init(); err != nil {
		slog.Error(err.Error(), "name", name, "offset", offset)
	}
	return p
}

// Shutdown resets the GPIO line allowing use by another program
func (gpio *GPIO) Shutdown() {
	for _, p := range gpio.Pins {
		p.Reconfigure(gpiocdev.AsInput)
		p.Close()
	}
	gpio.Pins = nil
}

// String returns the string representation of the GPIO
func (gpio *GPIO) String() string {
	str := ""
	for _, pin := range gpio.Pins {
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
