package devices

import (
	"fmt"

	"periph.io/x/conn/v3/analog"
	"periph.io/x/devices/v3/ads1x15"
)

// AnalogPin is the analog equivalent of a digital gpio pin
type AnalogPin struct {
	name string
	ads1x15.PinADC
}

func (p *AnalogPin) Name() string {
	return p.name
}

func (p *AnalogPin) String() string {
	return p.name + ": todo write String() function"
}

// Get returns a single float64 reading from the pin
func (p AnalogPin) Get() (float64, error) {
	reading, err := p.Read()
	if err != nil {
		return 0.0, err
	}
	var val float64
	fmt.Sscanf(reading.V.String(), "%f", &val)
	return val, err
}

// ReadContinous returns a channel that will continually read
// data from respective ads1115 pin and make the float64 values
// available as soon as the data is ready.
func (p AnalogPin) ReadContinuous() <-chan float64 {
	// Read values continuously from ADC.
	c := p.PinADC.ReadContinuous()
	floatQ := make(chan float64)
	go func() {
		var val float64
		var vsamp analog.Sample

		for {
			vsamp = <-c
			_, err := fmt.Sscanf(vsamp.V.String(), "%f", &val)
			if err != nil {
				panic(err)
			}
			floatQ <- val
		}
	}()

	return floatQ
}

// Close the pin and set it back to it's defaults. TODO
// set the pin back to its defaults
func (a AnalogPin) Close() {
	a.Halt()
}
