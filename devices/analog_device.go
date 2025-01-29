package devices

type AnalogDevice interface {
	Device
	Name() string
	ReadContinuous() <-chan float64
}

// NewAnalogDevice creates a analog device with the given name and ads1115
// pin value.  The options (not yet supported) should be those used by the ads1115
// library provided by
func NewAnalogDevice(name string, offset int, opts any) (ad AnalogDevice) {

	if Mock {
		return &AnalogMock{}
	}

	dev := &AnalogADS1115{}
	dev.BaseDevice = NewDevice(name)

	var err error
	// fix this only need to be done once, not for every pin
	ads1115 := GetADS1115()
	dev.AnalogPin, err = ads1115.Pin(name, offset, opts)
	if err != nil {
		panic(err)
	}
	return dev
}
