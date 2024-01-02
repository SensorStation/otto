package main

import (
	"log"

	"github.com/sensorstation/otto"
)

type controller struct {
}

// Globals
var (
	control    Control
	Controller controller
)

func (c controller) Init() error {
	o := otto.O()
	o.Subscribe("ss/d/#", Controller)
	return nil
}

func (c controller) Callback(t string, p []byte) {

	log.Printf("MQTT [I] Data: %s, %s", t, string(p))
	msg, err := otto.MsgFromMQTT(t, p)
	if err != nil {
		log.Printf("ERROR - parsing incoming message: %+v\n", err)
		return
	}

	control.Update(&msg.Data)

	// update the station that sent the msg
	// otto.Store.Store(msg)
	station := otto.Stations.Update(msg)
	if station == nil {
		log.Printf("Failed to update station for %+v\n", msg)
		return
	}

	o := otto.O()
	msg.Type = "station"
	o.InQ <- msg
}
