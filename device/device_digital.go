package device

import (
	"log/slog"

	"github.com/warthog618/go-gpiocdev"
)

// GPIODevice satisfies the Device interface and the GPIO Pin interface
type DigitalDevice struct {
	*Device
	*DigitalPin
}

// NewDigitalDevice.go creates a GPIO based device associated with GPIO Pin idx with the
// corresponding gpiocdev.LineReqOptions. If there is an gpiocdev.EventHandler an
// Event Q channel will be created.  Also a publication topic will be associated.
// For input devices the Pubs will be used to publish output state to the device.
// For output devices the Pubs will be used to publish the latest data collected
func NewDigitalDevice(name string, idx int, opts ...gpiocdev.LineReqOption) *DigitalDevice {
	d := &DigitalDevice{
		Device: NewDevice(name),
	}

	// look for an eventhandler, if so setup an event channel
	// for _, opt := range opts {
	// 	t := fmt.Sprintf("%T", opt)
	// 	if t == "gpiocdev.EventHandler" {
	// 		d.EvtQ = make(chan gpiocdev.LineEvent)
	// 	}
	// }

	// append the pubs
	gpio = GetGPIO()
	d.DigitalPin = gpio.Pin(name, idx, opts...)

	return d
}

func (d *DigitalDevice) ReadPub() {
	val, err := d.Get()
	if err != nil {
		slog.Error("ReadPub get value from pin", "device", d.Name, "error", err)
		return
	}
	d.Publish(val)
}

// import "github.com/warthog618/go-gpiocdev"

// func NewDigitalDevice(name string, offset int, opts ...gpiocdev.LineReqOption) *Device {
// 	d := NewDevice(name)
// 	g := GetGPIO()
// 	p := g.Pin(name, offset, opts...)
// 	d.ReadWriteCloser = p
// 	return d
// }

// type DigitalMock struct {
// }

// func GetDigitalMock(name string) *DigitalMock {
// 	return &DigitalMock{}
// }

// func (d *DigitalMock) Get() (int, error) {
// 	return 1, nil
// }

// func (d *DigitalMock) Set(val int) error {
// 	return nil
// }

// func (d *DigitalMock) On() error {
// 	return d.Set(1)
// }

// func (d *DigitalMock) Off() error {
// 	return d.Set(0)
// }

// func (d *DigitalMock) Close() error {
// 	return nil
// }

// func (d *DigitalMock) Read([]byte) (n int, err error) {
// 	return n, err
// }

// func (d *DigitalMock) Write([]byte) (n int, err error) {
// 	return n, err
// }
