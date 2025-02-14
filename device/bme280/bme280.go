package bme280

import (
	"encoding/json"
	"errors"

	"github.com/maciej/bme280"
	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/device/drivers"
)

// XXX - Maybe move this over to periph.io implementation?
// BME280 is an i2c device that gathers air temp, humidity and pressure
type BME280 struct {
	*device.Device

	bus    string
	addr   int
	driver *bme280.Driver
}

type Response bme280.Response

func New(name, bus string, addr int) *BME280 {
	b := &BME280{
		Device: device.NewDevice(name),
		bus:    bus,
		addr:   addr,
	}
	return b
}

// Init opens the i2c bus at the specified address and gets the device
// read for reading
func (b *BME280) Init() error {
	if device.IsMock() == true {
		return nil
	}

	i2c, err := drivers.GetI2CDriver(b.bus, b.addr)
	if err != nil {
		return err
	}

	b.driver = bme280.New(i2c)
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
	if device.IsMock() {
		return &bme280.Response{
			Temperature: 20.33,
			Pressure:    1027.33,
			Humidity:    74.33,
		}, nil
	}

	response, err := b.driver.Read()
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (b *BME280) ReadPub() error {
	vals, err := b.Read()
	if err != nil {
		return errors.New("Failed to read bme280: " + err.Error())
	}

	jb, err := json.Marshal(vals)
	if err != nil {
		return errors.New("BME280 failed marshal read response" + err.Error())
	}
	b.Publish(jb)
	return nil
}
