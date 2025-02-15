package drivers

import (
	"errors"
	"fmt"
	"log"

	"github.com/sensorstation/otto/device"
	"periph.io/x/conn/v3/analog"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

// ADS1115 is an i2c ADC chip that will use the i2c device type
// to provide 4 single analog pins to be used by the raspberry
// pi to access analog sensors via the i2c bus.  In a sense this
// device is a higher level device than the device_i2c.
type ADS1115 struct {
	pins [4]*ADS1115Pin
	bus  i2c.BusCloser
	adc  *ads1x15.Dev
	mock bool
}

var (
	ads1115 *ADS1115
)

// GetADS1115 will return the default ads1115 struct singleton. The
// first time GetADS1115 is called it will create a new device.
// Subsequent calls will return the global variable.
func GetADS1115() *ADS1115 {
	if ads1115 == nil {
		ads1115 = NewADS1115("ads1115", "/dev/i2c-1", 0x48)
	}
	return ads1115
}

// NewADS creates a new ADS1115 giving it the provided name,
// I2C bus (default /dev/i2c-1) and address (default 0x48).
func NewADS1115(name string, bus string, addr int) *ADS1115 {
	a := &ADS1115{}
	if device.IsMock() {
		a.mock = true
		return a
	}

	a.Init()
	return a
}

// Init prepares the chip for usage
func (a *ADS1115) Init() (err error) {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Printf("device_ads1115: host init failed: %s", err)
		return err
	}

	// Open default I²C bus.
	a.bus, err = i2creg.Open("")
	if err != nil {
		log.Printf("device_ads1115: i2c open failed: %s", err)
		return err
	}

	// Create a new ADS1115 ADC.
	a.adc, err = ads1x15.NewADS1115(a.bus, &ads1x15.DefaultOpts)
	if err != nil {
		log.Printf("device_ads1115: new ads failed: %s", err)
		return err
	}
	return nil
}

// Pin allocates and prepares one of the ads1115 pins (0 - 3) for use.
func (a *ADS1115) Pin(name string, ch int, opts any) (pin *ADS1115Pin, err error) {
	// Obtain an analog pin from the ADC.
	if ch < 0 || ch > 3 {
		return pin, fmt.Errorf("PinInit Invalid channel %d", ch)
	}

	c := ads1x15.Channel0
	switch ch {
	case 1:
		c = ads1x15.Channel1

	case 2:
		c = ads1x15.Channel2

	case 3:
		c = ads1x15.Channel3
	}

	pin = &ADS1115Pin{}
	pin.PinADC, err = a.adc.PinForChannel(c, 3300*physic.MilliVolt, 1*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		return pin, err
	}
	return pin, err
}

// Close the ads1115 and shutdown all the pins
func (a *ADS1115) Close() error {
	a.bus.Close()
	for i := 0; i < 4; i++ {
		if a.pins[i] != nil {
			a.pins[i].Close()
		}
	}
	return nil
}

func (a *ADS1115) String() string {
	return "ADS1115 todo write string function"
}

func (a *ADS1115) JSON() []byte {
	panic("write ads1115 JSON function")
	return nil
}

// ADS1115Pin is the analog equivalent of a digital gpio pin
type ADS1115Pin struct {
	name string
	ads1x15.PinADC
}

func (p *ADS1115Pin) Name() string {
	return p.name
}

func (p *ADS1115Pin) String() string {
	return p.name + ": todo write String() function"
}

// Get returns a single float64 reading from the pin
func (p ADS1115Pin) Read() (float64, error) {
	reading, err := p.PinADC.Read()
	if err != nil {
		return 0.0, err
	}
	var val float64
	// fmt.Sscanf(reading.V.String(), "%f", &val)
	val = float64(reading.V)
	return val, err
}

func (p ADS1115Pin) Set(val float64) error {
	return errors.New("Analog Pin ads1115 can not be set")
}

// ReadContinous returns a channel that will continually read
// data from respective ads1115 pin and make the float64 values
// available as soon as the data is ready.
func (p ADS1115Pin) ReadContinuous() <-chan float64 {
	// Read values continuously from ADC.
	c := p.PinADC.ReadContinuous()
	floatQ := make(chan float64)
	go func() {
		var val float64
		var vsamp analog.Sample

		for {
			vsamp = <-c
			_, err := fmt.Sscanf(vsamp.V.String(), "%f", &val)
			if err != nil {
				panic(err)
			}
			floatQ <- val
		}
	}()

	return floatQ
}

// Close the pin and set it back to it's defaults. TODO
// set the pin back to its defaults
func (a ADS1115Pin) Close() error {
	return a.Halt()
}
