/*
Relay sets up pin 6 for a Relay (or LED) and connects to an MQTT
broker waiting for instructions to turn on or off the relay.
*/

package main

import (
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/gpio"
	"github.com/warthog618/go-gpiocdev"
)

//go:embed app
var content embed.FS

type relay struct {
	*gpio.Pin
}

func main() {
	l := otto.GetLogger()

	// var data any
	s := otto.GetServer()
	// s.EmbedTempl("/", data, content)
	s.Appdir("/", "app")
	go s.Start()

	// Get the GPIO driver
	g := gpio.GetGPIO()
	defer func() {
		g.Shutdown()
	}()

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	r := g.Pin("relay", 6, gpiocdev.AsOutput(0))
	m := otto.GetMQTT()
	m.Connect()
	m.Subscribe("ss/c/station/relay", r)

	<-quit
	g.Shutdown()
	l.Info("Exiting relay")
}
