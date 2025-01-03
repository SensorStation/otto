package devices

type AnalogDevice struct {
	*Device
}

type AnalogPin struct {
	offset int
	val    float64
	mock   bool
}

func NewAnalogDevice(name string) *AnalogDevice {
	return &AnalogDevice{
		Device: NewDevice(name),
	}
}

func (a *AnalogPin) Get() float64 {
	return 0.0
}

func (a *AnalogPin) Set(v float64) {

}
