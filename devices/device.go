package devices

import (
	"fmt"
	"time"

	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/message"
	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

type Mode int

const (
	ModeNone Mode = iota
	ModeInput
	ModeOutput
	ModePWM
)

type Device struct {
	Name string
	Pub  string
	Subs []string
	Mode

	Period time.Duration
	EvtQ   chan gpiocdev.LineEvent
}

func NewDevice(name string) *Device {
	d := &Device{
		Name: name,
	}
	return d
}

func (d *Device) AddPub(p string) {
	d.Pub = p
}

func (d *Device) Subscribe(topic string, f func(*message.Msg)) {
	d.Subs = append(d.Subs, topic)
	m := messanger.GetMQTT()
	m.Subscribe(topic, f)
}

func (d *Device) Publish(data any) {
	if d.Pub == "" {
		logger.GetLogger().Error("Device.Publish failed has no pub", "name", d.Name)
		return
	}
	var buf []byte
	switch data.(type) {
	case []byte:
		buf = data.([]byte)

	case string:
		buf = []byte(data.(string))

	default:
		panic("unknown type: " + fmt.Sprintf("%T", data))
	}

	msg := message.New(d.Pub, buf, d.Name)
	m := messanger.GetMQTT()
	m.PublishMsg(msg)
}

func (d *Device) Shutdown() {
	GetGPIO().Shutdown()
}

func (d *Device) EventLoop(done chan bool, readpub func() error) {
	l := logger.GetLogger()

	running := true
	for running {
		select {
		case evt := <-d.EvtQ:
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

			l.Info("GPIO edge", "device", d.Name, "direction", evtype,
				"seqno", evt.Seqno, "lineseq", evt.LineSeqno)

			err := readpub()
			if err != nil {
				l.Error("Failed to read and publish", "device", d.Name, "error", err)
			}

		case <-done:
			running = false
		}
	}
}

func (d *Device) TimerLoop(done chan bool, readpub func() error) {
	// No need to loop if we don't have a ticker period
	if d.Period <= 0 {
		return
	}
	ticker := time.NewTicker(d.Period)

	running := true
	for running {
		select {
		case <-ticker.C:
			err := readpub()
			if err != nil {
				logger.GetLogger().Error("TimerLoop failed to readpub", "device", d.Name, "error", err)
			}

		case <-done:
			running = false
		}
	}
}
