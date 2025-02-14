package button

import (
	"log/slog"
	"time"

	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/device/drivers"
	"github.com/warthog618/go-gpiocdev"
)

type Button struct {
	*device.Device
	*drivers.DigitalPin
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
		Device:     device.NewDevice(name),
		DigitalPin: drivers.NewDigitalPin(name, offset, bopts...),
	}
	b.EvtQ = evtQ
	return b
}

func (b *Button) ReadPub() {
	val, err := b.Get()
	if err != nil {
		slog.Error("Failed to read buttons value: ", "error", err.Error())
		return
	}
	slog.Debug("read", "device", "button", "val", val)
	// val := strconv.Itoa(v)
	b.Publish(val)
}
