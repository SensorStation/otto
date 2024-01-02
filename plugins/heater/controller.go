package main

import "github.com/sensorstation/otto"

const (
	None = iota
	On
	Off
)

type Control struct {
	Max   float64
	Min   float64
	State int
}

func (c Control) On(station string) {
	otto.O().Publish("ss/c/"+station+"/heater", "on")
}

func (c Control) Off(station string) {
	otto.O().Publish("ss/c/"+station+"/heater", "off")
}

func (c Control) Update(msg *otto.MsgStation) {

	var v float64
	if v <= c.Min && c.State == Off {
		c.On(msg.ID)
	} else if v >= c.Max && c.State == On {
		c.Off(msg.ID)
	}
}
