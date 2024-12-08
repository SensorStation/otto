/*
Relay sets up pin 6 for a Relay (or LED) and connects to an MQTT
broker waiting for instructions to turn on or off the relay.
*/

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/gpio"
)

var (
	l *otto.Logger
)

type relay struct {
	*gpio.Pin
}

func main() {
	l = otto.GetLogger()

	// Get the GPIO driver
	g := gpio.GetGPIO()
	defer func() {
		g.Shutdown()
	}()

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	r := g.Pin("relay", 6, gpio.Output(0))
	m := otto.GetMQTT()
	m.Connect()
	m.Subscribe("/ss/d/station/relay", r)

	<-quit
	g.Shutdown()
	l.Info("Exiting relay")
}
