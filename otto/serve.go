package main

import (
	"github.com/sensorstation/otto"
	"github.com/spf13/cobra"
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

	// Start MQTT
	mqtt := otto.GetMQTT()
	mqtt.Connect()

	// start web server / rest server
	server = otto.NewServer()
	server.Start()

	// start websockets
	<-done
}
