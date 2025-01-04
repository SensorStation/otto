package devices

import "github.com/sensorstation/otto/utils"

type AnalogDevice struct {
	*Device
	*AnalogPin
}

type AnalogPin struct {
	name   string
	offset int
	val    float64
	mock   bool
}

func NewAnalogDevice(name string, pin int) *AnalogDevice {
	return &AnalogDevice{
		Device:    NewDevice(name),
		AnalogPin: NewAnalogPin(name, pin),
	}
}

func NewAnalogPin(name string, offset int) *AnalogPin {
	return &AnalogPin{
		name:   name,
		offset: offset,
	}
}

func (a *AnalogPin) Get() float64 {
	r := utils.NewRando()
	return r.Float64() * 100
}

func (a *AnalogPin) Set(v float64) {
	a.val = v
}
