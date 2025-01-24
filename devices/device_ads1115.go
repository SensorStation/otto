package devices

import (
	"fmt"

	"periph.io/x/conn/v3/analog"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

type ADS1115 struct {
	*I2CDevice
	Mock bool

	bus  i2c.BusCloser
	adc  *ads1x15.Dev
	pins [4]*AnalogPin
}

var (
	ads1115 *ADS1115
)

func GetADS1115() *ADS1115 {
	if ads1115 == nil {
		ads1115 = NewADS1115("ads1115", "/dev/i2c-1", 0x48)
		ads1115.Init()
	}
	return ads1115
}

func NewADS1115(name string, bus string, addr int) *ADS1115 {
	a := &ADS1115{
		I2CDevice: NewI2CDevice(name, bus, addr),
	}
	return a
}

func (a *ADS1115) Init() (err error) {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		return err
	}

	// Open default I²C bus.
	a.bus, err = i2creg.Open("")
	if err != nil {
		return err
	}

	// Create a new ADS1115 ADC.
	a.adc, err = ads1x15.NewADS1115(a.bus, &ads1x15.DefaultOpts)
	if err != nil {
		return err
	}
	return nil
}

func (a *ADS1115) Pin(name string, ch int, opts any) (pin AnalogPin, err error) {
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

	pin = AnalogPin{}
	pin.PinADC, err = a.adc.PinForChannel(c, 3300*physic.MilliVolt, 1*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		return pin, err
	}
	return pin, err
}

func (a *ADS1115) Close() {
	a.bus.Close()
	for i := 0; i < 4; i++ {
		if a.pins[i] != nil {
			a.pins[i].Halt()
		}
	}
}

type AnalogPin struct {
	ads1x15.PinADC
	lastread analog.Sample
}

func (p AnalogPin) Get() (float64, error) {
	reading, err := p.Read()
	if err != nil {
		return 0.0, err
	}
	var val float64
	p.lastread = reading
	fmt.Sscanf(reading.V.String(), "%f", &val)
	return val, err
}

func (p AnalogPin) ReadContinuous() <-chan float64 {
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
			p.lastread = vsamp
		}
	}()

	return floatQ
}

func (a AnalogPin) Close() {
	a.Halt()
}

// func (a *ADS1115) Read() (analog.Sample, error) {

// 	// Read values from ADC.
// 	reading, err := a.pin.Read()
// 	if err != nil {
// 		return reading, err
// 	}
// 	return reading, err
// }

// func (a *ADS1115) ReadPub() error {
// 	readQ := a.pin.ReadContinuous()
// 	for vals := range readQ {
// 		jb, err := json.Marshal(vals)
// 		if err != nil {
// 			return errors.New("BME280 failed marshal read response" + err.Error())
// 		}
// 		a.Publish(jb)
// 	}
// 	return nil
// }

// func (a *ADS1115) Halt() {
// 	a.pin.Halt()
// }
