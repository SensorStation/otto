package gpio

import "github.com/warthog618/go-gpiocdev"

// Line interface is used to emulate a GPIO pin as
// implemented by the go-gpiocdev package
type Line interface {
	Value() (int, error)
	SetValue(int) error
	Reconfigure(...gpiocdev.LineConfigOption) error
	Close() error
}

// MockGPIO fakes the Line interface on computers that don't
// actually have GPIO pins mostly for mocking tests
type MockLine struct {
	Val int
}

func (m MockLine) Value() (int, error) {
	return m.Val, nil
}

func (m MockLine) SetValue(val int) error {
	m.Val = val
	return nil
}

func (m MockLine) Reconfigure(...gpiocdev.LineConfigOption) error {
	return nil
}

func (m MockLine) Close() error {
	return nil
}
