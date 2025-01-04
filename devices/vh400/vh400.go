package vh400

import (
	"fmt"

	"github.com/sensorstation/otto/devices"
)

type VH400 struct {
	*devices.AnalogDevice
}

func New(name string, pin int) *VH400 {
	return &VH400{
		AnalogDevice: devices.NewAnalogDevice(name, pin),
	}
}

func (v *VH400) ReadPub() error {
	val := v.Get()
	vstr := fmt.Sprintf("%5.4f", val)
	v.Publish([]byte(vstr))
	return nil
}
