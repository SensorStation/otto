package main

import (
	"github.com/sensorstation/otto/gpio"
)

var (
	led = LED(6)
)

func main() {
	m := otto.GetMQTT()
	m.Start(broker)
	m.Subscribe("ss/c/localhost/led", led)
}

type LED struct {
	pin	int
}

func LED(pin int) {
	return &LED{
		pin: pin,
	}
}

