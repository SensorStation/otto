package main

import (
	"strconv"
	"time"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/gpio"
	"github.com/warthog618/go-gpiocdev"
)

var (
	l    *otto.Logger
	mqtt *otto.MQTT
)

func main() {

	l = otto.GetLogger()

	mqtt = otto.GetMQTT()
	mqtt.Connect()

	// Get the GPIO driver
	g := gpio.GetGPIO()
	defer func() {
		g.Shutdown()
	}()

	done := make(chan bool, 0)
	go startSwitchHandler(g, done)
	go startSwitchToggler(g, done)

	<-done
}

func startSwitchToggler(g *gpio.GPIO, done chan bool) {
	on := false
	r := g.Pin("reader", 23, gpiocdev.AsOutput(1))
	for {
		if on {
			r.On()
			on = false
		} else {
			r.Off()
			on = true
		}
		time.Sleep(1 * time.Second)
	}
}

func startSwitchHandler(g *gpio.GPIO, done chan bool) {
	evtQ := make(chan gpiocdev.LineEvent)
	sw := g.Pin("switch", 24, gpiocdev.WithPullUp, gpiocdev.WithBothEdges, gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
		evtQ <- evt
	}))

	for {
		select {
		case evt := <-evtQ:
			switch evt.Type {
			case gpiocdev.LineEventFallingEdge:
				l.Info("GPIO failing edge", "pin", sw.Name)
				fallthrough

			case gpiocdev.LineEventRisingEdge:
				l.Info("GPIO raising edge", "pin", sw.Name)

				v, err := sw.Get()
				if err != nil {
					otto.GetLogger().Error("Error getting input value: ", "error", err.Error())
					continue
				}
				val := strconv.Itoa(v)
				mqtt.Publish("ss/d/station/"+sw.Name, val)

			default:
				l.Warn("Unknown event type ", "type", evt.Type)
			}

		case <-done:
			return
		}
	}
}
