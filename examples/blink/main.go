/*
Blink sets up pin 6 for an LED and goes into an endless
toggle mode.
*/

package main

import (
	"fmt"
	"time"

	"github.com/sensorstation/otto/devices"
	"github.com/warthog618/go-gpiocdev"
)

func main() {

	// Get the GPIO driver
	g := devices.GetGPIO()
	defer func() {
		g.Shutdown()
	}()

	led := g.Pin("led", 6, gpiocdev.AsOutput(0))
	for {
		led.Toggle()
		fmt.Printf(g.String())
		time.Sleep(1 * time.Second)
	}
}
