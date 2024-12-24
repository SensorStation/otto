package main

import "github.com/sensorstation/otto/message"

type Relay struct {
	*GPIODevice
}

func NewRelay(name string, offset int) *Relay {
	relay := &Relay{
		GPIODevice: GPIOOut(name, offset),
	}
	return relay
}

func (r *Relay) Callback(msg *message.Msg) {
	switch msg.Path[3] {
	case "off":
		r.Off()

	case "on":
		r.On()
	}
}
