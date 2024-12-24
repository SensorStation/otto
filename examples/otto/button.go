package main

import (
	"fmt"
	"strconv"

	"github.com/sensorstation/otto"
	"github.com/warthog618/go-gpiocdev"
)

type Button struct {
	*GPIODevice
}

func NewButton(name string, pin int) *Button {
	b := &Button{
		GPIODevice: GPIOIn(name, pin),
	}
	b.pubs = append(b.pubs, "ss/c/"+stationName+"/"+name)
	return b
}

func (b *Button) ButtonLoop(done chan bool) {
	running := true
	for running {
		fmt.Printf("Button[%s:%d] waiting for eventQ %p\n", b.Name(), b.Offset(), b.evtQ)
		select {
		case evt := <-b.evtQ:
			evtype := "falling"
			switch evt.Type {
			case gpiocdev.LineEventFallingEdge:
				evtype = "falling"

			case gpiocdev.LineEventRisingEdge:
				evtype = "raising"

			default:
				l.Warn("Unknown event type ", "type", evt.Type)
				continue
			}

			l.Info("GPIO edge", "device", b.Name(), "direction", evtype,
				"seqno", evt.Seqno, "lineseq", evt.LineSeqno)

			v, err := b.Get()
			if err != nil {
				l.Error("Failed to read buttons value", "error", err)
				continue
			}

			val := strconv.Itoa(v)
			for _, t := range b.pubs {
				otto.GetMQTT().Publish(t, val)
			}

		case <-done:
			running = false
		}
	}
}
