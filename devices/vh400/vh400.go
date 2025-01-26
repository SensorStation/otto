package vh400

// See https://vegetronix.com/Products/VH400/VH400-Piecewise-Curve
// For calculations on the VWC.  Borrowed from above website

import (
	"log"

	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/messanger"
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
	volts := v.Get()
	vwc := v.vwc(volts)
	if v.Pub != "" {
		v.Publish(vwc)
	}

	return vwc
}

func (v *VH400) ReadPub() error {
	v.Pub = messanger.TopicData("vh100/" + v.Name)
	q := v.AnalogDevice.ReadContinuous()
	go func() {
		for {
			volts := <-q
			vwc := v.vwc(volts)
			v.Publish(vwc)
		}
	}()

	return nil
}

/*
Most curves can be approximated with linear segments of the form:

y= m*x-b,

where m is the slope of the line

The VH400's Voltage to VWC curve can be approximated with 4 segents of the form:

VWC= m*V-b

where V is voltage.

m = (VWC2 - VWC1) / (V2-V1)

where V1 and V2 are voltages recorded at the respective VWC levels
of VWC1 and VWC2. After m is determined, the y-axis intercept
coefficient b can be found by inserting one of the end points into
the equation:

b= m*v-VWC
*/
func (v *VH400) vwc(volts float64) float64 {
	var coef float64
	var rem float64

	switch {
	case volts >= 0.0 && volts <= 1.1:
		coef = 10.0
		rem = 1.0

	case volts > 1.1 && volts <= 1.3:
		coef = 25.0
		rem = 17.5

	case volts > 1.3 && volts <= 1.82:
		coef = 48.08
		rem = 47.5

	case volts > 1.82 && volts <= 2.2:
		coef = 26.32
		rem = 7.80

	case volts > 2.2 && volts <= 3.0:
		coef = 62.5
		rem = 7.89

	default:
		log.Printf("Invalide voltage out of range 0.0 -> 3.0 %5.2f", volts)
		return 0.0
	}
	vwc := coef*volts - rem
	return vwc
}
