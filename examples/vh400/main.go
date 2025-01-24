/*
Relay sets up pin 6 for a Relay (or LED) and connects to an MQTT
broker waiting for instructions to turn on or off the relay.
*/

package main

import (
	"fmt"

	"github.com/sensorstation/otto/devices/vh400"
)

func main() {
	soil := vh400.New("vh400", 0)
	for v := range soil.ReadContinuous() {
		fmt.Printf("%+v\n", v)
	}
}
