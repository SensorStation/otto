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

	r := &relay{g.Pin("relay", 6, gpio.Output(0))}
	m := otto.GetMQTT()
	m.Connect()
	m.Subscribe("/ss/d/station/relay", r)

	l.Info("Recieved signal exiting")
	<-quit
	leave()
}

func (r *relay) SubCallback(topic string, data []byte) {
	msg := otto.NewMsg(topic, data)

	switch msg.String() {
	case "on":
		r.On()

	case "off":
		r.Off()

	case "toggle":
		r.Toggle()

	case "exit":
		l.Info("recieved exit command")
		leave()

	default:
		l.Warn("relay unknown command", "msg", msg.String())
	}
}

func leave() {
	l.Info("Exiting relay")
	os.Exit(0)
}
