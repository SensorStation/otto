package main

import "github.com/sensorstation/otto"

var (
	mqtt     *otto.MQTT
	server   *otto.Server
	stations *otto.StationManager
)

func main() {
	Execute()
}
