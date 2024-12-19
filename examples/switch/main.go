package main

import (
	"fmt"
	"time"

	"github.com/sensorstation/otto/gpio"
	"github.com/warthog618/go-gpiocdev"
)

var in *gpio.Pin

func main() {

	// Get the GPIO driver
	g := gpio.GetGPIO()
	defer func() {
		g.Shutdown()
	}()

	on := false

	sw := g.Pin("switch", 23, gpiocdev.AsOutput(1))
	in = g.Pin("in", 24, gpiocdev.WithPullUp, gpiocdev.WithBothEdges, gpiocdev.WithEventHandler(eventHandler))
	for {
		if on {
			sw.On()
			on = false
		} else {
			sw.Off()
			on = true
		}

		v, err := in.Get()
		fmt.Printf("%d - error: %v\n", v, err)
		time.Sleep(1 * time.Second)
	}
}

func eventHandler(evt gpiocdev.LineEvent) {
	fmt.Println("we had an event")
	t := time.Now()
	edge := "rising"
	if evt.Type == gpiocdev.LineEventFallingEdge {
		edge = "falling"
	}
	if evt.Seqno != 0 {
		// only uAPI v2 populates the sequence numbers
		fmt.Printf("event: #%d(%d)%3d %-7s %s (%s)\n",
			evt.Seqno,
			evt.LineSeqno,
			evt.Offset,
			edge,
			t.Format(time.RFC3339Nano),
			evt.Timestamp)
	} else {
		fmt.Printf("event:%3d %-7s %s (%s)\n",
			evt.Offset,
			edge,
			t.Format(time.RFC3339Nano),
			evt.Timestamp)
	}
	v, err := in.Get()
	fmt.Printf("%d - error: %v\n", v, err)
}
