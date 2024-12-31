package main

import "github.com/sensorstation/otto/devices"

type SSD1306 struct {
	*devices.I2CDevice
}

func New(name string, bus, addr int) {

}
