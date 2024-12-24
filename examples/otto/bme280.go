package main

import (
	"fmt"

	"github.com/maciej/bme280"
	"golang.org/x/exp/io/i2c"
)

// BME280 is an i2c device that gathers air temp, humidity and pressure
type BME280 struct {
	Addr   int
	Bus    string
	driver *bme280.Driver
	name   string
}

func NewBME280(name, bus string, addr int) *BME280 {
	b := &BME280{
		name: name,
		Addr: addr,
		Bus:  bus,
	}
	return b
}

// Init opens the i2c bus at the specified address and gets the device
// read for reading
func (b *BME280) Init() error {
	device, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, b.Addr)
	if err != nil {
		return err
	}

	b.driver = bme280.New(device)
	err = b.driver.InitWith(bme280.ModeForced, bme280.Settings{
		Filter:                  bme280.FilterOff,
		Standby:                 bme280.StandByTime1000ms,
		PressureOversampling:    bme280.Oversampling16x,
		TemperatureOversampling: bme280.Oversampling16x,
		HumidityOversampling:    bme280.Oversampling16x,
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BME280) Name() string {
	return b.name
}

func (b *BME280) Read() (*bme280.Response, error) {

	response, err := b.driver.Read()
	if err != nil {
		return nil, err
	}

	fmt.Printf("response: %+v\n", response)
	return &response, err
}
