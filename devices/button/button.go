package button

import (
	"errors"
	"strconv"
	"time"

	"github.com/sensorstation/otto/devices"
	"github.com/warthog618/go-gpiocdev"
)

type Button struct {
	*devices.GPIODevice
}

func New(name string, pin int) *Button {
	b := &Button{}
	b.GPIODevice = devices.NewGPIODevice(name, pin, devices.ModeInput,
		gpiocdev.WithPullUp,
		gpiocdev.WithFallingEdge,
		gpiocdev.WithDebounce(10*time.Millisecond),
		gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
			b.EvtQ <- evt
		}))
	return b
}

func (b *Button) ReadPub() error {
	v, err := b.Get()
	if err != nil {
		return errors.New("Failed to read buttons value: " + err.Error())
	}

	val := strconv.Itoa(v)
	b.Publish(val)
	return nil
}
