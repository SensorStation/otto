package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type PinWR struct {
	gpio.PinIO
}

func GetPinRW(name string) (p PinWR) {
	if p.PinIO = gpioreg.ByName(name); p.PinIO == nil {
		log.Fatalln("Could not find pin: ", name)
	}

	if err := p.In(gpio.Float, gpio.NoEdge); err != nil {
		log.Fatalln("Could not set input params on pin", name)
	}
	return p
}

func (p PinWR) Get() interface{} {
	var o bool
	l := p.PinIO.Read()
	if l == gpio.Low {
		o = false
	} else {
		o = true
	}
	return o
}

func (p PinWR) GetLevel() gpio.Level {
	return p.PinIO.Read()
}

func (p PinWR) Set(i interface{}) {
	d := i.(gpio.Level)
	p.Out(d)
}

func (p PinWR) SetLevel(l gpio.Level) {
	p.PinIO.Out(l)
}

func (p PinWR) MessageHandler(c mqtt.Client, m mqtt.Message) {
	log.Printf("Received message on topic: %s\nMessage: %s\n", m.Topic(), m.Payload())
}

type FakePin struct {
	Path string
}

func (p FakePin) MsgHandler(c mqtt.Client, m mqtt.Message) {
	if config.Debug {
		log.Printf("Received message on topic: %s\nMessage: %s\n", m.Topic(), m.Payload())
	}
	payload := m.Payload()
	fmt.Printf("Rando Payload %+v\n", string(payload))
}
