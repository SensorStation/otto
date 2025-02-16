package drivers

import (
	"math/rand"
	"time"
)

type AnalogPin interface {
	Read() (float64, error)
	ReadContinuous() <-chan float64
	Close() error
}

type MockAnalogPin struct {
	PinName string
	Offset  int
	val     float64
}

func NewMockAnalogPin(name string, pin int, opts ...any) *MockAnalogPin {
	return &MockAnalogPin{
		PinName: name,
		Offset:  pin,
	}
}

func (a *MockAnalogPin) Read() (float64, error) {
	a.val = rand.Float64()
	return a.val, nil
}

func (a *MockAnalogPin) ReadContinuous() <-chan float64 {
	readQ := make(chan float64)

	go func() {
		for {
			val, _ := a.Read()
			readQ <- val
			<-time.After(1 * time.Second)
		}
	}()

	return readQ
}

func (a *MockAnalogPin) Close() error {
	return nil
}
