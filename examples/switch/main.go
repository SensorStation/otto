package main

import (
	"time"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/gpio"
	"github.com/warthog618/go-gpiocdev"
)

var (
	l *otto.Logger
)

func main() {

	l = otto.GetLogger()

	m := otto.GetMQTT()
	m.Connect()

	// Get the GPIO driver
	g := gpio.GetGPIO()
	defer func() {
		g.Shutdown()
	}()

	done := make(chan bool, 0)
	startEventHandler(g, done)
	startButtonToggler(g, done)
}

func startButtonToggler(g *gpio.GPIO, done chan bool) {
	on := false
	sw := g.Pin("switch", 23, gpiocdev.AsOutput(1))
	for {
		if on {
			sw.On()
			on = false
		} else {
			sw.Off()
			on = true
		}
		time.Sleep(1 * time.Second)
	}
}

func startEventHandler(g *gpio.GPIO, done chan bool) {
	evtQ := make(chan gpiocdev.LineEvent)
	in := g.Pin("in", 24, gpiocdev.WithPullUp, gpiocdev.WithBothEdges, gpiocdev.WithEventHandler(func (evt gpiocdev.LineEvent) {
		evtQ <- evt
	}))

	for {
		select {
		case evt := <-evtQ:
			switch evt.Type {
			case gpiocdev.LineEventFallingEdge:
				l.Info("GPIO failing edge", "pin", in.Name)
				fallthrough

			case gpiocdev.LineEventRisingEdge:
				l.Info("GPIO raising edge", "pin", in.Name)

				v, err := in.Get()
				if err != nil {
					otto.GetLogger().Error("Error getting input value: ", "error", err.Error())
					continue
				}
				mqtt := otto.GetMQTT()
				mqtt.Publish("ss/c/station/" + in.Name, v)

			default:
				l.Warn("Unknown event type ", "type", evt.Type)
			}

		case <-done:
			return
		}
	}
}


// func eventHandler(evt gpiocdev.LineEvent) {
// 	t := time.Now()
// 	edge := "rising"
// 	if evt.Type == gpiocdev.LineEventFallingEdge {
// 		edge = "falling"
// 	}
// 	if evt.Seqno != 0 {
// 		// only uAPI v2 populates the sequence numbers
// 		fmt.Printf("event: #%d(%d)%3d %-7s %s (%s)\n",
// 			evt.Seqno,
// 			evt.LineSeqno,
// 			evt.Offset,
// 			edge,
// 			t.Format(time.RFC3339Nano),
// 			evt.Timestamp)
// 	} else {
// 		fmt.Printf("event:%3d %-7s %s (%s)\n",
// 			evt.Offset,
// 			edge,
// 			t.Format(time.RFC3339Nano),
// 			evt.Timestamp)
// 	}
// 	v, err := in.Get()
// }
