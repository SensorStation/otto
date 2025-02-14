package devices

import (
	"time"

	"periph.io/x/conn/v3/analog"
	"periph.io/x/devices/v3/ads1x15"
)

type AnalogMock struct {
	BaseDevice
}

func (a *AnalogMock) Get() float64 {
	return 54.5
}

func (a *AnalogMock) ReadContinuous() <-chan float64 {
	q := make(chan float64)

	go func() {
		for {
			v := a.Get()
			q <- v
			time.Sleep(1 * time.Second)
		}
	}()

	return q
}

type MockAnalogPin struct {
	ads1x15.PinADC
}

func (ma *MockAnalogPin) Name() string {
	return ""
}

func (ma *MockAnalogPin) Number() int {
	return 0
}

func (ma *MockAnalogPin) ReadContinuous() (ch <-chan analog.Sample) {

	return ch
}

func (ma *MockAnalogPin) Range() (s1 analog.Sample, s2 analog.Sample) {
	return s1, s2
}

func (ma *MockAnalogPin) Read() (s analog.Sample, err error) {
	return s, err
}

func (ma *MockAnalogPin) Function() string {
	return ""
}

func (ma *MockAnalogPin) String() string {
	return ""
}

func (ma *MockAnalogPin) Halt() error {
	return nil
}
