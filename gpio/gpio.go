package main

import (
	"fmt"

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
	Description string `json:"description"`
	Offset      int    `json:"offset"`
	Name        string `json:"name"`

	Mode `json:"mode"`
	*gpiocdev.Line
}

var (
	gpio *GPIO
)

func (pin *Pin) String() string {
	str := fmt.Sprintf("%s : %s - pin %4d", pin.Name, pin.Description, pin.Offset)
	return str
}

func Output(val int) Mode {
	return ModeOutput
}

func (pin *Pin) Get() (int, error) {
	return pin.Line.Value()
}

func (pin *Pin) Set(v int) error {
	return pin.Line.SetValue(v)
}

func (pin *Pin) On() error {
	return pin.Set(1)
}

func (pin *Pin) Off() error {
	return pin.Set(0)
}

type GPIO struct {
	*gpiocdev.Chip
	chipname string
	Pins     map[int]*Pin `json:"pins"`
}

func GetGPIO() *GPIO {
	gpio := &GPIO{
		chipname: "gpiochip4",
	}
	gpio.Pins = make(map[int]*Pin)
	return gpio
}

func (gpio *GPIO) Pin(desc string, offset int, mode Mode) (p *Pin) {

	l, err := gpiocdev.RequestLine(gpio.chipname, offset, gpiocdev.AsOutput(0))
	if err != nil {
		panic(err)
	}

	p = &Pin{
		Description: desc,
		Offset:      offset,
		Line:        l,
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

// func (gpio *GPIO) Find(offset int) (p *Pin) {

// 	l, err := gpiocdev.RequestLine(gpio.chipname, offset, gpiocdev.AsInput)
// 	if err != nil {
// 		fmt.Printf("Finding line %s returned error: %s\n", gpio.chipname, err)
// 		os.Exit(1)
// 	}

// 	p = &Pin{
// 		Line:   l,
// 		Offset: offset,
// 	}
// 	return p
// }
