package main

import (
	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/gpio"
)

var (
	l    *otto.Logger
	done chan bool
)

func main() {
	l = otto.GetLogger()
	done = make(chan bool)

	// TODO capture signals

	initSignals()
	initApp()
	initStations()
	devices.initDevices(done)
	// cmd.Execute()

	<-done
	cleanup()
}

func cleanup() {
	g := gpio.GetGPIO()
	g.Shutdown()
}

func initSignals() {
	// todo make sure we capture signals
}

func initStations() {
	st := otto.GetStationManager()
	st.Start()
}
