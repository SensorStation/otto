package gpio

import (
	"fmt"

	"github.com/sensorstation/otto"
	"github.com/warthog618/go-gpiocdev"
)

type Mode int

const (
	ModeNone Mode = iota
	ModeOutput
	ModeInput
	ModePWM
)

type Pin struct {
	Name   string `json:"name"`
	Offset int    `json:"offset"`
	Value  int    `json:"value"`

	Mode `json:"mode"`

	*gpiocdev.Line
}

var (
	gpio *GPIO
)

func (pin *Pin) String() string {
	str := fmt.Sprintf("%s - pin %4d", pin.Name, pin.Offset)
	return str
}

func Output(val int) Mode {
	return ModeOutput
}

func (pin *Pin) Get() (int, error) {
	return pin.Line.Value()
}

func (pin *Pin) Set(v int) error {
	pin.Value = v
	return pin.Line.SetValue(v)
}

func (pin *Pin) On() error {
	return pin.Set(1)
}

func (pin *Pin) Off() error {
	return pin.Set(0)
}

func (pin *Pin) Toggle() error {
	val := ^pin.Value
	return pin.Set(val)
}

func (pin Pin) SubCallback(t string, payload []byte) {
	val := string(payload)
	switch val {
	case "on":
		pin.On()

	case "off":
		pin.Off()

	case "toggle":
		pin.Toggle()

	}
}

func (p *Pin) Input() error {
	return p.Reconfigure(gpiocdev.AsInput)
}

func (p *Pin) Output(v int) error {
	return p.Reconfigure(gpiocdev.AsOutput(v))
}

type GPIO struct {
	*gpiocdev.Chip
	chipname string
	Pins     map[int]*Pin `json:"pins"`
}

func GetGPIO() *GPIO {
	if gpio == nil {
		gpio = &GPIO{
			chipname: "gpiochip4", // raspberry pi
		}
		gpio.Pins = make(map[int]*Pin)
	}
	return gpio
}

func (gpio *GPIO) Pin(name string, offset int, mode Mode) (p *Pin) {
	l := otto.GetLogger()

	// Put this into emulation mode if it fails
	line, err := gpiocdev.RequestLine(gpio.chipname, offset, gpiocdev.AsOutput(0))
	if err != nil {
		l.Error("Failed to get gpio line", "error", err)
		return nil
	}

	p = &Pin{
		Name:   name,
		Offset: offset,
		Line:   line,
	}

	gpio.Pins[offset] = p
	return p
}

func (gpio *GPIO) Shutdown() {
	for _, p := range gpio.Pins {
		p.Reconfigure(gpiocdev.AsInput)
		p.Close()
	}
}

func (gpio *GPIO) String() string {
	str := ""
	for p, pin := range gpio.Pins {
		str += fmt.Sprintf("%s pin %d\n", pin.String(), p)
	}
	return str
}
