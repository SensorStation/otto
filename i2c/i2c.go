package i2c

import (
	"machine"
)

type I2C struct {
	Device string
	_ = machine.I2C{}
}
