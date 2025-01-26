package devices

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

type DeviceIntf interface {
}

type Mode int

const (
	ModeNone Mode = iota
	ModeInput
	ModeOutput
	ModePWM
)

// Device is an abstract
type Device struct {
	// Name of the device human readable
	Name string
	// Suffix to be appended to the base topic for mqtt publications
	Pub string
	// Subscription this device will listen to
	Subs []string
	// Input, Output, PWM, I2C, SPI, UART, etc
	Mode

	// Period for repititive timed tasks like collecting and
	// publishing data
	Period time.Duration

	// EventQ for devices that are interupt driven
	EvtQ chan gpiocdev.LineEvent
}

// NewDevice creates a new device with the given name
func NewDevice(name string) *Device {
	d := &Device{
		Name: name,
	}
	return d
}

// AddPub adds a publication
func (d *Device) AddPub(p string) {
	d.Pub = p
}

func (d *Device) Subscribe(topic string, f func(*messanger.Msg)) {
	d.Subs = append(d.Subs, topic)
	m := messanger.GetMQTT()
	m.Subscribe(topic, f)
}

func (d *Device) Publish(data any) {
	if d.Pub == "" {
		slog.Error("Device.Publish failed has no pub", "name", d.Name)
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

	msg := messanger.New(d.Pub, buf, d.Name)
	m := messanger.GetMQTT()
	m.PublishMsg(msg)
}

func (d *Device) Shutdown() {
	// XXX = this is not right
	// GetGPIO().Shutdown()
}

func (d *Device) EventLoop(done chan bool, readpub func() error) {
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
				slog.Warn("Unknown event type ", "type", evt.Type)
				continue
			}

			slog.Info("GPIO edge", "device", d.Name, "direction", evtype,
				"seqno", evt.Seqno, "lineseq", evt.LineSeqno)

			err := readpub()
			if err != nil {
				slog.Error("Failed to read and publish", "device", d.Name, "error", err)
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
				slog.Error("TimerLoop failed to readpub", "device", d.Name, "error", err)
			}

		case <-done:
			running = false
		}
	}
}
