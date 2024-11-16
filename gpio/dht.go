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
	WAIT_DURATION  = 2 * time.Microsecond
	DHT_DATA_BITS  = 40
	DHT_DATA_BYTES = (DHT_DATA_BITS / 8)
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
	var p, i int
	d.Input()

	start := time.Now()
	ts := time.Since(start)
	for ts < timeout {
		// time.Sleep(1 * time.Microsecond)
		ts = time.Since(start)

		p, err = d.Get()
		if err == nil && p == expectedState {
			return ts, nil
		}
		i++
	}
	fmt.Printf("timeout[%d]: %v - %v - expected %d got %d\n",
		i, timeout, ts, expectedState, p)
	return ts, fmt.Errorf("Timeout")
}

func (d *DHT) Read() (err error) {

	// early allocations before time critical code
	var data []byte
	data = make([]byte, DHT_DATA_BYTES)

	if err = d.Input(); err != nil {
		return err
	}

	on, err := d.Get()
	if err != nil {
		fmt.Println("could not read the thing")
		return err
	}
	fmt.Printf("0. Are we ready? %d\n", on)

	// fmt.Println("A. MCU pulls low for 1-10ms - we wait 8ms")
	// MCU Pulls down for 1-10 ms // dht.cc 18 ms
	if err = d.Output(0); err != nil {
		return err
	}
	time.Sleep(8 * time.Millisecond)

	// fmt.Println("B. Pull high and wait 20 - 40 us")
	// MCU will pull the signal high and wait 20-40us
	if err = d.On(); err != nil {
		return err
	}

	if _, err := d.waitPinState(40*time.Microsecond, 0); err != nil {
		fmt.Printf("phase B %v\n", err)
		return err
	}

	// fmt.Println("C. wait 88ms for the pin to be pulled high")
	if _, err := d.waitPinState(88*time.Microsecond, 1); err != nil {
		fmt.Printf("phase C %v\n", err)
		return err
	}

	// fmt.Println("D. wait 88ms for the pin to be polled low")
	if _, err := d.waitPinState(88*time.Microsecond, 0); err != nil {
		fmt.Printf("phase E %v\n", err)
		return err
	}

	// fmt.Println("E. Time to read 40 bits of data")
	for i := 0; i < DHT_DATA_BITS; i++ {
		lowdur, err := d.waitPinState(65*time.Microsecond, 1)
		if err != nil {
			return fmt.Errorf("Reading low duration: %s", err)
		}

		highdur, err := d.waitPinState(75*time.Microsecond, 0)
		if err != nil {
			return fmt.Errorf("Reading high duration: %s", err)
		}

		b := i / 8
		m := i % 8
		if m == 0 {
			data[b] = 0
		}
		if highdur > lowdur {
			data[b] |= 1 << (7 - m)
		} else {
			data[b] |= 0 << (7 - m)
		}
		fmt.Printf("%d: high: %v - low: %v\n", i, highdur, lowdur)
	}

	err = d.Output(1)
	if err != nil {
		fmt.Printf("Error setting output to one")
	}

	for i := 0; i < len(data); i++ {
		fmt.Printf("0x%02x ", data[i])
	}

	// fmt.Println("")
	// if data[4] != ((data[0] + data[1] + data[2] + data[3]) & 0xFF) {
	// 	// return fmt.Errorf("Data Checksum failed")
	// }

	// var val uint16
	// val = uint16(data[0]) & 0x7f
	// val <<= 8
	// val |= uint16(data[1])
	// if data[0] > 127 {
	// 	val = -val
	// }
	// d.humidity = (float32)(val / 10)
	// fmt.Printf("h: %d %5.2f\n", val, d.humidity)

	// val = uint16(data[2]) & 0x7f
	// val <<= 8
	// val |= uint16(data[3])
	// if data[2] > 127 {
	// 	val = -val
	// }
	// fmt.Printf("t: %d %5.2f\n", val, d.temperature)
	// d.temperature = (float32)(val / 10)

	// convert to bytes
	// data := make([]uint8, 5)
	for i := range data {
		for j := 0; j < 8; j++ {
			data[i] <<= 1
			if lengths[i*8+j] > LOGICAL_1_TRESHOLD {
				data[i] |= 0x01
			}
		}
	}

	if err := d.checksum(data); err != nil {
		if err != nil {
			return err
		}
	}

	var (
		humidity    uint16
		temperature uint16
	)

	// calculate humidity

	humidity |= uint16(data[0])
	humidity <<= 8
	humidity |= uint16(data[1])

	if humidity < 0 || humidity > 1000 {
		return HumidityError
	}

	d.humidity = float32(humidity) / 10

	// calculate temperature
	temperature |= uint16(data[2])
	temperature <<= 8
	temperature |= uint16(data[3])

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

func (d *DHT) checksum(data []uint8) error {
	var sum uint8

	for i := 0; i < 4; i++ {
		sum += data[i]
	}

	if sum != data[4] {
		return ChecksumError
	}

	return nil
}
