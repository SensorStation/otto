package gpio

import (
	"encoding/json"
	"fmt"

	"github.com/sensorstation/otto"
	"github.com/warthog618/go-gpiocdev"
)

// Mode causes the Pin to be configured as Input, Output,
// PWM or Analog
type Mode int

// Mode values
const (
	ModeNone Mode = iota
	ModeOutput
	ModeInput
	ModePWM
	ModeAnalog
)

type Line interface {
}

// Pin represents a single GPIO Pin
type Pin struct {
	Name   string `json:"name"`
	Offset int    `json:"offset"`
	Value  int    `json:"value"`
	Mode   `json:"mode"`

	*gpiocdev.Line
}

var (
	gpio *GPIO
)

// Init the pin from the offset and mode
func (p *Pin) Init() error {

	gpio := GetGPIO()

	// Put this into emulation mode if it fails
	// line, err := gpiocdev.RequestLine(gpio.Chipname, p.Offset, p.Mode)
	line, err := gpiocdev.RequestLine(gpio.Chipname, p.Offset, gpiocdev.AsOutput(0))
	if err != nil {
		return err
	}
	p.Line = line
	return nil
}

// String returns a string representation of the GPIO pin
func (pin *Pin) String() string {
	str := fmt.Sprintf("%s - pin %4d", pin.Name, pin.Offset)
	return str
}

// Output sets the pin to be an output with the default value
func Output(val int) Mode {
	return ModeOutput
}

// Get returns the value of the pin, an error is returned if
// the GPIO value fails
func (pin *Pin) Get() (int, error) {
	if pin.Line == nil {
		return 0, fmt.Errorf("GPIO not active")
	}
	return pin.Line.Value()
}

// Set the value of the pin
func (pin *Pin) Set(v int) error {
	pin.Value = v
	if pin.Line == nil {
		return fmt.Errorf("GPIO not active")
	}
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
	val := ^pin.Value
	return pin.Set(val)
}

// SubCallback is the default callback for pins if they are
// registered with the MQTT.Subscribe() function
func (pin Pin) SubCallback(t string, d []byte) {
	msg := otto.NewMsg(t, d, "mqtt-pin-"+pin.Name)

	fmt.Printf("mqtt msg: %+v", msg)

	switch msg.String() {
	case "on":
		pin.On()

	case "off":
		pin.Off()

	case "toggle":
		pin.Toggle()

	}
}

// Input configures the given pin as an pinput
func (p *Pin) Input() error {
	return p.Reconfigure(gpiocdev.AsInput)
}

// Output configures the given pin as an output
func (p *Pin) Output(v int) error {
	return p.Reconfigure(gpiocdev.AsOutput(v))
}

// GPIO is used to initialize the GPIO and pins on a raspberry pi
type GPIO struct {
	Chipname string       `json:"chipname"`
	Pins     map[int]*Pin `json:"pins"`
}

// GetGPIO returns the GPIO singleton for the Raspberry PI
func GetGPIO() *GPIO {
	if gpio == nil {
		gpio = &GPIO{
			Chipname: "gpiochip4", // raspberry pi
		}
		gpio.Pins = make(map[int]*Pin)
	}
	return gpio
}

// Init initialized the GPIO
func (gpio *GPIO) Init() error {
	l := otto.GetLogger()
	for _, pin := range gpio.Pins {
		if err := pin.Init(); err != nil {
			l.Error("Error initializing pin ", "name", pin.Name, "offset", pin.Offset)
		}
	}
	return nil
}

// Pin initializes the given GPIO pin, name and mode
func (gpio *GPIO) Pin(name string, offset int, mode Mode) (p *Pin) {
	l := otto.GetLogger()

	p = &Pin{
		Name:   name,
		Offset: offset,
	}

	gpio.Pins[offset] = p
	if err := p.Init(); err != nil {
		l.Error(err.Error(), "name", name, "offset", offset)
	}
	return p
}

// Shutdown resets the GPIO line allowing use by another program
func (gpio *GPIO) Shutdown() {
	for _, p := range gpio.Pins {
		p.Reconfigure(gpiocdev.AsInput)
		p.Close()
	}
}

// String returns the string representation of the GPIO
func (gpio *GPIO) String() string {
	str := ""
	for p, pin := range gpio.Pins {
		str += fmt.Sprintf("%s pin %d\n", pin.String(), p)
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
