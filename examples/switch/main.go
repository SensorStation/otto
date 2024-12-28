package main

import (
	"strconv"
	"time"

	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/messanger"
	"github.com/warthog618/go-gpiocdev"
)

var (
	l    *logger.Logger
	mqtt *messanger.MQTT
)

func main() {

	l = logger.GetLogger()

	mqtt = messanger.GetMQTT()

	// Get the GPIO driver
	g := devices.GetGPIO()
	defer func() {
		g.Shutdown()
	}()

	done := make(chan bool, 0)
	go startSwitchHandler(g, done)
	go startSwitchToggler(g, done)

	<-done
}

func startSwitchToggler(g *devices.GPIO, done chan bool) {
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

func startSwitchHandler(g *devices.GPIO, done chan bool) {
	evtQ := make(chan gpiocdev.LineEvent)
	sw := g.Pin("switch", 24, gpiocdev.WithPullUp, gpiocdev.WithBothEdges, gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
		evtQ <- evt
	}))

	for {
		select {
		case evt := <-evtQ:
			switch evt.Type {
			case gpiocdev.LineEventFallingEdge:
				l.Info("GPIO failing edge", "pin", sw.Offset())
				fallthrough

			case gpiocdev.LineEventRisingEdge:
				l.Info("GPIO raising edge", "pin", sw.Offset())
				v, err := sw.Get()
				if err != nil {
					logger.GetLogger().Error("Error getting input value: ", "error", err.Error())
					continue
				}
				val := strconv.Itoa(v)
				mqtt.Publish("ss/d/station/switch", val)

			default:
				l.Warn("Unknown event type ", "type", evt.Type)
			}

		case <-done:
			return
		}
	}
}
