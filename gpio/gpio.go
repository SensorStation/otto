package main

import (
	"fmt"
	"log"

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

func (pin Pin) Callback(t string, payload []byte) {
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
			chipname: "gpiochip4",
		}
		gpio.Pins = make(map[int]*Pin)
	}
	return gpio
}

func (gpio *GPIO) Pin(name string, offset int, mode Mode) (p *Pin) {

	fmt.Printf("GPIO: %v\n", gpio)

	l, err := gpiocdev.RequestLine(gpio.chipname, offset, gpiocdev.AsOutput(0))
	if err != nil {
		panic(err)
	}

	p = &Pin{
		Name:   name,
		Offset: offset,
		Line:   l,
	}

	gpio.Pins[offset] = p

	if mqtt != nil {
		log.Printf("mode: %d\n", mode)

		switch mode {
		case ModeInput:
			log.Println("mode input")

		case ModeOutput:
			log.Println("mode output")
			mqtt.Subscribe("ss/station/dev/"+p.Name, p)

		}
	}

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
