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
	j, err := g.JSON()
	if err != nil {
		fmt.Printf("Failed to JSONify GPIO: %v\n", err)
	}
	fmt.Printf("GPIO: %s", j)
	for {
		led.Toggle()
		time.Sleep(1 * time.Second)
	}
}
