package main

import (
	"time"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/gpio"
	"github.com/warthog618/go-gpiocdev"
)

type Device interface {
	Name() string
}

type GPIODevice struct {
	*gpio.Pin

	subs []string
	pubs []string

	evtQ chan gpiocdev.LineEvent
}

func (d *GPIODevice) Name() string {
	return d.Pin.Name
}

func (d *GPIODevice) On() {
	d.Set(1)
}

func (d *GPIODevice) Off() {
	d.Set(0)
}

func (d *GPIODevice) Set(v int) {
	d.Pin.Set(v)
	val := "off"
	if v > 0 {
		val = "on"
	}

	m := otto.GetMQTT()
	for _, p := range d.pubs {
		m.Publish(p, val)
	}
}

func GPIOOut(name string, pin int) *GPIODevice {
	g := gpio.GetGPIO()
	d := &GPIODevice{
		Pin: g.Pin(name, pin, gpiocdev.AsOutput(0)),
	}

	topic := "ss/c/" + stationName + "/" + name
	d.subs = append(d.subs, topic)

	m := otto.GetMQTT()
	m.Subscribe(topic, d)

	d.pubs = append(d.pubs, "ss/d/"+stationName+"/"+name)
	return d
}

func GPIOIn(name string, pin int) *GPIODevice {
	g := gpio.GetGPIO()
	d := &GPIODevice{
		evtQ: make(chan gpiocdev.LineEvent),
	}
	d.Pin = g.Pin(name, pin,
		gpiocdev.WithPullUp,
		gpiocdev.WithFallingEdge,
		gpiocdev.WithDebounce(10*time.Millisecond),
		gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
			d.evtQ <- evt
		}))

	d.pubs = append(d.pubs, "ss/c/"+stationName+"/"+name)
	return d
}

func (dm *DeviceManager) initDevices(done chan bool) {
	m := otto.GetMQTT()

	relay := NewRelay("relay", 22)
	m.Subscribe("ss/c/"+stationName+"/off", relay)
	m.Subscribe("ss/c/"+stationName+"/on", relay)
	devices.Add(relay)

	led := NewLED("led", 6)
	m.Subscribe("ss/c/"+stationName+"/off", led)
	m.Subscribe("ss/c/"+stationName+"/on", led)
	devices.Add(led)

	onButton := NewButton("on", 23)
	devices.Add(onButton)

	offButton := NewButton("off", 27)
	devices.Add(offButton)

	bme := NewBME280("bme", "/dev/i2c-1", 0x77)
	devices.Add(bme)

	go onButton.ButtonLoop(done)
	go offButton.ButtonLoop(done)
}
