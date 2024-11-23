package main

import (
	"github.com/sensorstation/otto"
	"github.com/spf13/cobra"
)

var (
	mqtt     *otto.MQTT
	server   *otto.Server
	stations *otto.StationManager
	sensors  *otto.SensorManager
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start oTTo the Server",
	Long:  `Start OttO the IoT Server`,
	Run:   serveRun,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serveRun(cmd *cobra.Command, args []string) {

	done := make(chan interface{})

	// Allocate and start the station manager
	stations = otto.NewStationManager()
	stations.Start()

	// Open the data storage
	sensors = otto.NewSensorManager()
	mqtt_init()

	// start web server / rest server
	server = otto.NewServer()
	server.Start()

	// start websockets
	<-done
}

func mqtt_init() {
	// Start MQTT
	mqtt := otto.GetMQTT()
	mqtt.Connect()

	// ss/d/<station>/<sensor>
	mqtt.Subscribe("ss/d/+/+", sensors)
}
