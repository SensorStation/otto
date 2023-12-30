package main

import "github.com/rustyeddy/iote"

const (
	None = iota
	On
	Off
)

type Controller struct {
	Max   float64
	Min   float64
	State int
}

func (c Controller) On(station string) {
	mqtt.Publish("ss/c/"+station+"/heater", "on")
}

func (c Controller) Off(station string) {
	mqtt.Publish("ss/c/"+station+"/heater", "off")
}

func (c Controller) Update(msg *iote.MsgStation) {

	var v float64
	if v <= c.Min && c.State == Off {
		c.On(msg.ID)
	} else if v >= c.Max && c.State == On {
		c.Off(msg.ID)
	}
}
