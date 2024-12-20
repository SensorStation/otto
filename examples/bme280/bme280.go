package main

import (
	"fmt"

	"github.com/sensorstation/otto/i2c/bme280"
	"golang.org/x/exp/io/i2c"
)

type BME280 struct {
	Addr   int
	Dev    string
	driver *bme280.Driver
}

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

func (b *BME280) Read() (*bme280.Response, error) {

	response, err := b.driver.Read()
	if err != nil {
		return nil, err
	}

	fmt.Printf("response: %+v\n", response)
	return &response, err
}
