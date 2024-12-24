package main

import (
	"fmt"
	"encoding/json"
	"time"

	"github.com/maciej/bme280"
	"github.com/sensorstation/otto"
	"golang.org/x/exp/io/i2c"
)

// BME280 is an i2c device that gathers air temp, humidity and pressure
type BME280 struct {
	Addr   int
	Bus    string
	driver *bme280.Driver
	name   string
	Period time.Duration
	Pubs   []string
}

func NewBME280(name, bus string, addr int) *BME280 {
	b := &BME280{
		name:   name,
		Addr:   addr,
		Bus:    bus,
		Period: 10 * time.Second,
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

	fmt.Printf("driver: %+v\n", b.driver)
	response, err := b.driver.Read()
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (b *BME280) Loop(done chan bool) {
	timer := time.NewTimer(b.Period)

	running := true
	for running {
		select {
		case <-timer.C:
			vals, err := b.Read()
			if err != nil {
				l.Error("Failed to read bme280", "error", err)
				continue
			}

			jb, err := json.Marshal(vals)
			if err != nil {
				otto.GetLogger().Error("failed to unmarshal bme Response", "error", err.Error())
				done <- true
				break
			}
			mqtt := otto.GetMQTT()
			for _, t := range b.Pubs {
				mqtt.Publish(t, jb)
			}

		case <-done:
			running = false
		}
	}
}
