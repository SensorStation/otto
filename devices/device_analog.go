package devices

// AnalogGPIO is a wrapper around the ads1115 to provide
// analog (float64) values
type AnalogGPIO interface {
	Pin(name string, offset int, opts any) *AnalogPin
}

// APin is the Analog version of a GPIO pin.
type APin interface {
	Get() (float64, error)
	ReadContinuous() <-chan float64
	Close()
}

// AnalogDevice is an instatiation of the GPIO device
type AnalogDevice struct {
	*Device
	AnalogPin
}

// NewAnalogDevice creates a analog device with the given name and ads1115
// pin value.  The options (not yet supported) should be those used by the ads1115
// library provided by
func NewAnalogDevice(name string, offset int, opts any) (ad *AnalogDevice) {
	// fix this only need to be done once, not for every pin
	ads1115 := GetADS1115()
	ad = &AnalogDevice{
		Device: NewDevice(name),
	}

	var err error
	ad.AnalogPin, err = ads1115.Pin(name, offset, opts)
	if err != nil {
		panic(err)
	}
	return ad
}
