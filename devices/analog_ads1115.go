package devices

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

// AnalogDevice is an instatiation of the GPIO device
type AnalogADS1115 struct {
	*BaseDevice
	*AnalogPin
}

// ADS1115 is an i2c ADC chip that will use the i2c device type
// to provide 4 single analog pins to be used by the raspberry
// pi to access analog sensors via the i2c bus.  In a sense this
// device is a higher level device than the device_i2c.
type ADS1115 struct {
	*I2CDevice
	Mock bool
	pins [4]*AnalogPin

	bus i2c.BusCloser
	adc *ads1x15.Dev
}

var (
	ads1115 *ADS1115
)

func (a AnalogADS1115) Name() string {
	return a.BaseDevice.Name()
}

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
	a := &ADS1115{
		I2CDevice: NewI2CDevice(name, bus, addr),
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
func (a *ADS1115) Pin(name string, ch int, opts any) (pin *AnalogPin, err error) {
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

	pin = &AnalogPin{}
	pin.PinADC, err = a.adc.PinForChannel(c, 3300*physic.MilliVolt, 1*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		return pin, err
	}
	return pin, err
}

// Close the ads1115 and shutdown all the pins
func (a *ADS1115) Close() {
	a.bus.Close()
	for i := 0; i < 4; i++ {
		if a.pins[i] != nil {
			a.pins[i].Close()
		}
	}
}

func (a *ADS1115) String() string {
	return "ADS1115 todo write string function"
}

func (a *ADS1115) JSON() []byte {
	panic("write ads1115 JSON function")
	return nil
}
