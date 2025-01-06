/*
Blink sets up pin 6 for an LED and goes into an endless
toggle mode.
*/

package main

import (
	"log/slog"
	"time"

	"encoding/json"

	"github.com/sensorstation/otto/devices"
)

var gpioStr = `
{
    "chipname":"gpiochip4",
    "pins": {
        "6": {
            "name": "led",
            "offset": 6,
            "value": 0,
            "mode": 0
        }
    }
}
`

func main() {
	var g devices.GPIO
	if err := json.Unmarshal([]byte(gpioStr), &g); err != nil {
		slog.Error(err.Error())
		return
	}

	if err := g.Init(); err != nil {
		slog.Error(err.Error())
		return
	}

	defer func() {
		g.Shutdown()
	}()

	led := g.Pins[6]
	for {
		led.Toggle()
		time.Sleep(1 * time.Second)
	}
}
