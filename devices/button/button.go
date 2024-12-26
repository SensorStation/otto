package button

import (
	"strconv"
	"time"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/devices"
	"github.com/warthog618/go-gpiocdev"
)

type Button struct {
	*devices.DeviceGPIO
}

func New(name string, pin int) *Button {
	b := &Button{}
	b.DeviceGPIO = devices.NewDeviceGPIO(name, pin, devices.ModeInput,
		gpiocdev.WithPullUp,
		gpiocdev.WithFallingEdge,
		gpiocdev.WithDebounce(10*time.Millisecond),
		gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
			b.EvtQ <- evt
		}))
	b.Pubs = append(b.Pubs, otto.TopicControl(name))
	return b
}

func (b *Button) EventLoop(done chan bool) {
	l := otto.GetLogger()

	running := true
	for running {
		select {
		case evt := <-b.EvtQ:
			evtype := "falling"
			switch evt.Type {
			case gpiocdev.LineEventFallingEdge:
				evtype = "falling"

			case gpiocdev.LineEventRisingEdge:
				evtype = "raising"

			default:
				l.Warn("Unknown event type ", "type", evt.Type)
				continue
			}

			l.Info("GPIO edge", "device", b.Name, "direction", evtype,
				"seqno", evt.Seqno, "lineseq", evt.LineSeqno)

			v, err := b.Get()
			if err != nil {
				l.Error("Failed to read buttons value", "error", err)
				continue
			}

			val := strconv.Itoa(v)
			for _, t := range b.Pubs {
				otto.GetMQTT().Publish(t, val)
			}

		case <-done:
			running = false
		}
	}
}
