package drivers

import (
	"fmt"

	"golang.org/x/exp/io/i2c"
)

var i2cbuses map[string]*i2cbus

func init() {
	i2cbuses = make(map[string]*i2cbus)
}

// represents a single i2c bus
type i2cbus struct {
	bus     string
	devices map[int]*i2c.Device
}

func GetI2CDriver(bus string, addr int) (device *i2c.Device, err error) {
	b := getI2CBus(bus)
	if b == nil {
		return device, fmt.Errorf("failed to get I2C bus %s", bus)
	}
	var ex bool
	if device, ex = b.devices[addr]; !ex {
		device, err = b.open(addr)
		if err != nil {
			return device, err
		}
		b.devices[addr] = device
	}
	return device, err
}

func getI2CBus(bus string) (b *i2cbus) {
	var ex bool
	if b, ex = i2cbuses[bus]; !ex {
		b = &i2cbus{
			bus:     bus,
			devices: make(map[int]*i2c.Device),
		}
		i2cbuses[bus] = b
	}
	return b
}

func (i *i2cbus) open(addr int) (dev *i2c.Device, err error) {
	d, err := i2c.Open(&i2c.Devfs{i.bus}, addr)
	if err != nil {
		return dev, err
	}
	i.devices[addr] = d
	return d, err
}
