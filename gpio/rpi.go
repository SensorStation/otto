package gpio

import (
	"fmt"
	"os"

	gpiocdev "github.com/warthog618/go-gpiocdev"
)

type Mode int

const (
	ModeNone Mode = iota
	ModeOutput
	ModeInput
	ModePWM
)

type Pin struct {
	Description string `json:description`
	Offset      int    `json:offset`
	Name        string `json:name`

	Mode `json:mode`
	*gpiocdev.Line
}

type RPI struct {
	Chip     *gpiocdev.Chip
	ChipName string       `json:chip-name`
	Pins     map[int]*Pin `json:pins`
}

var (
	rpi *RPI
)

func GetRPI() *RPI {
	if rpi == nil {
		rpi = &RPI{
			ChipName: "gpiochip4",
		}
	}
	rpi.Pins = make(map[int]*Pin)
	return rpi
}

func (rpi *RPI) Output(val int) Mode {
	return ModeOutput
}

func (pin *Pin) String() string {
	str := fmt.Sprintf("%4d: %s", pin.Offset, pin.Description)
	return str
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

func (rpi *RPI) Find(offset int) (p *Pin) {

	l, err := gpiocdev.RequestLine(rpi.ChipName, offset, gpiocdev.AsInput)
	if err != nil {
		fmt.Printf("Finding line %s returned error: %s\n", rpi.ChipName, err)
		os.Exit(1)
	}

	p = &Pin{
		Line:   l,
		Offset: offset,
	}
	return p
}

func (rpi *RPI) PinInit(desc string, offset int, mode Mode) (p *Pin) {

	// m := gpiocdev.LineDirectionInput
	// if mode == ModeOutput {
	// 	m = gpiocdev.LineDirectionOutput
	// }

	l, err := gpiocdev.RequestLine(rpi.ChipName, offset, gpiocdev.AsOutput(0))
	if err != nil {
		panic(err)
	}

	p = &Pin{
		Description: desc,
		Offset:      offset,
		Line:        l,
	}

	rpi.Pins[offset] = p
	return p
}

func (rpi *RPI) String() string {
	str := ""
	for p, pin := range rpi.Pins {
		str += fmt.Sprintf("%4d: %s\n", p, pin.String())
	}
	return str
}
