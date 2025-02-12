package device

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sensorstation/otto/messanger"
)

type Opener interface {
	Open() error
}

type OnOff interface {
	On()
	Off()
}

var mock bool

func Mock(mocking bool) {
	mock = mocking
}

func IsMock() bool {
	return mock
}

type Name interface {
	Name() string
}

type Device struct {
	// Name of the device human readable
	name string

	// Suffix to be appended to the base topic for mqtt publications
	pub string

	// Subscription this device will listen to
	subs []string

	// Period for repititive timed tasks like collecting and
	// publishing data
	period time.Duration

	// EventQ for devices that are interupt driven
	// EvtQ chan gpiocdev.LineEvent

	// for mocking
	val any

	// Last Error encountered
	Error error

	Opener
}

// NewDevice creates a new device with the given name
func NewDevice(name string) *Device {
	d := &Device{
		name: name,
	}
	return d
}

func (d Device) Name() string {
	return d.name
}

// AddPub adds a publication
func (d *Device) AddPub(p string) {
	d.pub = p
}

func (d *Device) GetPub() string {
	return d.pub
}

// func (d *Device) EventQ() (evtQ chan<- gpiocdev.LineEvent) {
// 	return d.EvtQ
// }

func (d *Device) Publish(data any) {
	if d.pub == "" {
		slog.Error("Device.Publish failed has no pub", "name", d.Name)
		return
	}
	var buf []byte

	m := messanger.GetMQTT()
	switch data.(type) {
	case []byte:
		buf = data.([]byte)

	case string:
		buf = []byte(data.(string))

	case int:
		m.Publish(d.pub, data)
		return

	default:
		panic("unknown type: " + fmt.Sprintf("%T", data))
	}

	msg := messanger.New(d.pub, buf, d.name)
	m.PublishMsg(msg)
}

func (d *Device) GetSubs() []string {
	return d.subs
}

func (d *Device) Subscribe(topic string, f func(*messanger.Msg)) {
	d.subs = append(d.subs, topic)
	m := messanger.GetMQTT()
	m.Subscribe(topic, f)
}

func (d *Device) TimerLoop(period time.Duration, done chan any, readpub func() error) {
	// No need to loop if we don't have a ticker period
	d.period = period
	if d.period <= 0 {
		return
	}
	ticker := time.NewTicker(d.period)

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

func (d *Device) String() string {
	return d.Name() + ": todo finish String()"
}

func (d *Device) JSON() []byte {
	panic("todo finish Device.JSON()")
}

func (d *Device) Close() error {
	// XXX = this is not right
	// GetGPIO().Shutdown()
	return nil
}
