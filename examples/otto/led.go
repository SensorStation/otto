package main

import "github.com/sensorstation/otto/message"

type LED struct {
	*GPIODevice
}

func NewLED(name string, offset int) *LED {
	led := &LED{
		GPIODevice: GPIOOut(name, offset),
	}
	return led
}

func (l *LED) Callback(msg *message.Msg) {
	switch msg.Path[3] {
	case "off":
		l.Off()

	case "on":
		l.On()
	}
}
