package devices

import (
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

type Pin interface {
	Name() string
}

type DigitalPin struct {
	name string
	opts []gpiocdev.LineReqOption
	Line

	offset int
	val    int
	mock   bool
}

// Init the pin from the offset and mode
func (p *DigitalPin) Init() error {

	gpio := GetGPIO()
	if gpio.Mock {
		line := GetMockLine(p.offset, p.opts...)
		p.mock = true
		p.Line = line
		return nil
	}

	line, err := gpiocdev.RequestLine(gpio.Chipname, p.offset, p.opts...)
	if err != nil {
		return err
	}
	p.Line = line
	return nil
}

func (p *DigitalPin) SetOpts(opts ...gpiocdev.LineReqOption) {
	p.opts = append(p.opts, opts...)
}

// String returns a string representation of the GPIO pin
func (pin *DigitalPin) String() string {
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
func (pin *DigitalPin) Get() (int, error) {
	if pin.Line == nil {
		return 0, fmt.Errorf("GPIO not active")
	}
	return pin.Line.Value()
}

// Set the value of the pin. Note: you can NOT set the value
// of an input pin, so we will check it and return an error.
// This maybe worthy of making it a panic!
func (pin *DigitalPin) Set(v int) error {
	if pin.Line == nil {
		return fmt.Errorf("GPIO not active")
	}
	pin.val = v
	return pin.Line.SetValue(v)
}

// On sets the value of the pin to 1
func (pin *DigitalPin) On() error {
	return pin.Set(1)
}

// Off sets the value of the pin to 0
func (pin *DigitalPin) Off() error {
	return pin.Set(0)
}

// Toggle with flip the value of the pin from 1 to 0 or 0 to 1
func (pin *DigitalPin) Toggle() error {
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
func (pin DigitalPin) Callback(msg *messanger.Msg) {
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
