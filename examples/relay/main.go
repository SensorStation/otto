/*
Relay sets up pin 6 for a Relay (or LED) and connects to an MQTT
broker waiting for instructions to turn on or off the relay.
*/

package main

import (
	"embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/messanger"
	"github.com/sensorstation/otto/server"
	"github.com/warthog618/go-gpiocdev"
)

//go:embed app
var content embed.FS

type relay struct {
	*devices.DigitalPin
}

func main() {
	// var data any
	s := server.GetServer()
	// s.EmbedTempl("/", data, content)
	s.Appdir("/", "app")
	go s.Start()

	// Get the GPIO driver
	g := devices.GetGPIO()
	defer func() {
		g.Close()
	}()

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	r := g.Pin("relay", 6, gpiocdev.AsOutput(0))
	m := messanger.GetMQTT()
	m.Connect()
	m.Subscribe("ss/c/station/relay", r.Callback)

	<-quit
	slog.Info("Exiting relay")
}
