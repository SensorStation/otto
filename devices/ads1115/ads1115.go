package ads1115

import (
	"fmt"

	"github.com/sensorstation/otto/devices"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

type ADS1115 struct {
	*devices.I2CDevice
	Bus  string
	Mock bool

	pin ads1x15.PinADC
}

func New(name string, bus string, addr int) *ADS1115 {
	a := &ADS1115{
		I2CDevice: devices.NewI2CDevice(name, bus, addr),
		Bus:       "/dev/i2c-1",
	}
	return a
}

func (a *ADS1115) Init() error {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		return err
	}

	// Open default I²C bus.
	// bus, err := i2creg.Open(a.Bus)
	bus, err := i2creg.Open("")
	if err != nil {
		return err
	}
	defer bus.Close()
	fmt.Printf("bus: %+v\n", bus)

	// Create a new ADS1115 ADC.
	adc, err := ads1x15.NewADS1115(bus, &ads1x15.DefaultOpts)
	if err != nil {
		return err
	}

	// Obtain an analog pin from the ADC.
	a.pin, err = adc.PinForChannel(ads1x15.Channel0, 3300*physic.MilliVolt, 1*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		return err
	}
	return nil
}

func (a *ADS1115) Read() error {

	// Read values from ADC.
	fmt.Println("Single reading")
	fmt.Printf("pin %+v\n", a.pin)
	// reading, err := a.pin.Read()
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("R: ", reading)

	// Read values continuously from ADC.
	fmt.Println("Continuous reading")
	c := a.pin.ReadContinuous()

	for reading := range c {
		fmt.Println("RC: ", reading)
	}
	return nil
}

func (a *ADS1115) ReadPub() error {
	// readQ, err := a.Read()
	// if err != nil {
	// 	return errors.New("Failed to read bme280: " + err.Error())
	// }

	// jb, err := json.Marshal(vals)
	// if err != nil {
	// 	return errors.New("BME280 failed marshal read response" + err.Error())
	// }
	// a.Publish(jb)
	// return nil
	return nil
}

func (a *ADS1115) Halt() {
	a.pin.Halt()
}

func (a *ADS1115) Close() {
	defer a.pin.Halt()
}
