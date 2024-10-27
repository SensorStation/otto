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

type Pinmap struct {
	Chip     *gpiocdev.Chip
	ChipName string       `json:chip-name`
	Pins     map[int]*Pin `json:pins`
}

var (
	pm Pinmap
)

func init() {
	pm.ChipName = "gpiochip4"
}

func (pin *Pin) Init(desc string, offset int, mode Mode) *Pin {
	p := pm.PinInit(desc, offset, mode)
	return p
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

func (pin *Pin) Off() {
	return pin.Set(0)
}

func (pm *Pinmap) Find(offset int) (p *Pin) {

	l, err := gpiocdev.RequestLine(pm.ChipName, offset, gpiocdev.AsInput)
	if err != nil {
		fmt.Printf("Finding line %s returned error: %s\n", pm.ChipName, err)
		os.Exit(1)
	}

	p = &Pin{
		Line:   l,
		Offset: offset,
	}
	return p
}

func (pm *Pinmap) PinInit(desc string, offset int, mode Mode) (p *Pin) {
	p = &Pin{
		Description: desc,
		Offset:      offset,
	}

	p.Init(desc, offset, mode)
	pm.Pins[offset] = p
	return p
}

func (pm *Pinmap) String() string {
	str := ""
	for p, pin := range pm.Pins {
		str += fmt.Sprintf("%4d: %s\n", p, pin.String())
	}
	return str
}
