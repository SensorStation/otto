package vh400

import (
	"github.com/sensorstation/otto/devices"
)

type VH400 struct {
	*devices.AnalogDevice // i2c device
}

func New(name string, pin int) *VH400 {
	v := &VH400{
		AnalogDevice: devices.NewAnalogDevice("vh400", pin, nil),
	}
	return v
}

func (v *VH400) Get() float64 {
	val := v.Get()

	// Todo some conversions

	return val
}

func (v *VH400) ReadPub() error {

	return nil
}
