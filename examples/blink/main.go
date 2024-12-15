/*
Blink sets up pin 6 for an LED and goes into an endless
toggle mode.
*/

package main

import (
	"fmt"
	"time"

	"github.com/sensorstation/otto/gpio"
)

func main() {

	// Get the GPIO driver
	g := gpio.GetGPIO()
	defer func() {
		g.Shutdown()
	}()

	led := g.Pin("led", 6, gpio.Output(0))
	for {
		led.Toggle()
		fmt.Printf(g.String())
		time.Sleep(1 * time.Second)
	}
}
