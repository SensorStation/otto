package main

import (
	"errors"
	"fmt"
	"time"
)

type DHT struct {
	*Pin

	temperature float32
	humidity    float32
	error
}

const (
	WAIT_DURATION = 2 * time.Microsecond
)

var (
	ChecksumError    = errors.New("checksum error")
	HumidityError    = errors.New("humidity range error")
	TemperatureError = errors.New("temperature range error")
)

func NewDHT(name string, offset int, mode Mode) *DHT {
	d := &DHT{}
	d.Pin = gpio.Pin(name, offset, mode)
	return d
}

func (d *DHT) Temperature() float32 {

	return 0.0
}

func (d *DHT) Humidity() float32 {

	return 0.0
}

func (d *DHT) waitPinState(timeout time.Duration, expectedState int) (duration time.Duration, err error) {

	d.Input()

	var elapsed time.Duration
	start := time.Now()

	for time.Since(start) < timeout {

		time.Sleep(2 * time.Microsecond)
		if p, err := d.Get(); err == nil {
			if p == expectedState {
				return elapsed, nil
			}
		}
	}
	return elapsed, fmt.Errorf("Timeout")
}

func (d *DHT) Read() (err error) {

	// early allocations before time critical code
	lengths := make([]time.Duration, 40)

	// MCU Pulls down for 1-10 ms
	if err = d.Output(0); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	// MCU will pull the signal high and wait 20-40us
	if err = d.On(); err != nil {
		return err
	}

	// phase B - 40us
	if _, err := d.waitPinState(40*time.Microsecond, 0); err != nil {
		return err
	}

	if _, err := d.waitPinState(88*time.Microsecond, 1); err != nil {
		return err
	}

	if _, err := d.waitPinState(88*time.Microsecond, 0); err != nil {
		return err
	}

	var totdur time.Duration
	for i := 0; i < 40; i++ {
		dur, err := d.waitPinState(65, 1)
		if err != nil {
			return err
		}
		totdur += dur

		dur, err = d.waitPinState(75, 0)
		if err != nil {
			return err
		}
		totdur += dur
	}

	// convert to bytes
	bytes := make([]uint8, 5)
	for i := range bytes {
		for j := 0; j < 8; j++ {
			bytes[i] <<= 1
			if lengths[i*8+j] > 10 /*LOGICAL_1_TRESHOLD*/ {
				bytes[i] |= 0x01
			}
		}
	}

	if err := d.checksum(bytes); err != nil {
		if err != nil {
			return err
		}
	}

	var (
		humidity    uint16
		temperature uint16
	)

	// calculate humidity

	humidity |= uint16(bytes[0])
	humidity <<= 8
	humidity |= uint16(bytes[1])

	if humidity < 0 || humidity > 1000 {
		return HumidityError
	}

	d.humidity = float32(humidity) / 10

	// calculate temperature
	temperature |= uint16(bytes[2])
	temperature <<= 8
	temperature |= uint16(bytes[3])

	// check for negative temperature
	if temperature&0x8000 > 0 {
		d.temperature = float32(temperature&0x7FFF) / -10
	} else {
		d.temperature = float32(temperature) / 10
	}

	// datasheet operating range
	if d.temperature < -40 || d.temperature > 80 {
		return TemperatureError
	}

	return err
}

func (d *DHT) checksum(bytes []uint8) error {
	var sum uint8

	for i := 0; i < 4; i++ {
		sum += bytes[i]
	}

	if sum != bytes[4] {
		return ChecksumError
	}

	return nil
}
