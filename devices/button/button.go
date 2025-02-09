package button

import (
	"errors"
	"log/slog"
	"time"

	"github.com/sensorstation/otto/device"
	"github.com/warthog618/go-gpiocdev"
)

type Button struct {
	*device.DigitalDevice
}

func New(name string, offset int, opts ...gpiocdev.LineReqOption) *Button {
	var evtQ chan gpiocdev.LineEvent
	evtQ = make(chan gpiocdev.LineEvent)
	bopts := []gpiocdev.LineReqOption{
		gpiocdev.WithPullUp,
		gpiocdev.WithDebounce(10 * time.Millisecond),
		gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
			evtQ <- evt
		}),
	}
	for _, o := range opts {
		bopts = append(bopts, o)
	}
	b := &Button{
		DigitalDevice: device.NewDigitalDevice(name, offset, bopts...),
	}
	b.EvtQ = evtQ
	return b
}

func (b *Button) Pub() error {
	var buf []byte
	n, err := b.Read(buf)
	if err != nil {
		return errors.New("Failed to read buttons value: " + err.Error())
	}
	slog.Debug("read", "device", "button", "bytes", n)
	// val := strconv.Itoa(v)
	b.Publish(buf)
	return nil
}
