package main

import "github.com/sensorstation/otto"

// all the global variables
var (
	mqtt *otto.MQTT
)

func main() {
	Execute()
}
