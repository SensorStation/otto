package devices

// AnalogDevice is an instatiation of the GPIO device
type AnalogDevice struct {
	BaseDevice
	*AnalogPin
}

// NewAnalogDevice creates a analog device with the given name and ads1115
// pin value.  The options (not yet supported) should be those used by the ads1115
// library provided by
func NewAnalogDevice(name string, offset int, opts any) (ad AnalogDevice) {
	// fix this only need to be done once, not for every pin
	ads1115 := GetADS1115()
	d := NewDevice(name)
	ad = AnalogDevice{
		BaseDevice: *d,
	}

	var err error
	ad.AnalogPin, err = ads1115.Pin(name, offset, opts)
	if err != nil {
		panic(err)
	}
	return ad
}

func (a AnalogDevice) Name() string {
	return a.BaseDevice.Name()
}
