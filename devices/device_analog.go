package devices

type AnalogGPIO interface {
	Pin(name string, offset int, opts any) *AnalogPin
}

type APin interface {
	Get() (float64, error)
	ReadContinuous() <-chan float64
	Close()
}

type AnalogDevice struct {
	*Device
	AnalogPin
}

func NewAnalogDevice(name string, offset int, opts any) (ad *AnalogDevice) {
	ads1115 := GetADS1115()
	ad = &AnalogDevice{
		Device: NewDevice("name"),
	}

	var err error
	ad.AnalogPin, err = ads1115.Pin("name", offset, opts)
	if err != nil {
		panic(err)
	}
	return ad
}

// type AnalogPin struct {
// 	name   string
// 	offset int
// 	val    float64
// 	mock   bool
// }

// func NewAnalogDevice(name string, pin int) *AnalogDevice {
// 	return &AnalogDevice{
// 		Device:    NewDevice(name),
// 		AnalogPin: NewAnalogPin(name, pin),
// 	}
// }

// func NewAnalogPin(name string, offset int) *AnalogPin {
// 	return &AnalogPin{
// 		name:   name,
// 		offset: offset,
// 	}
// }

// func (a *AnalogPin) Get() float64 {
// 	r := utils.NewRando()
// 	return r.Float64() * 100
// }

// func (a *AnalogPin) Set(v float64) {
// 	a.val = v
// }
