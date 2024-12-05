package main

import (
	"fmt"

	"github.com/sensorstation/otto"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start oTTo the Server",
	Long:  `Start OttO the IoT Server`,
	Run:   serveRun,
}

var (
	foreground bool
	done       chan bool
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().BoolVar(&foreground, "foreground", false, "Run the server command in the foreground")
}

func serveRun(cmd *cobra.Command, args []string) {

	if done != nil {
		// server has already been started
		fmt.Println("Server has already been started")
		return
	}
	done = make(chan bool)

	// Allocate and start the station manager
	stations := otto.GetStationManager()
	stations.Start()

	mqtt := otto.GetMQTT()
	err := mqtt.Connect()
	if err != nil {
		l.Println("MQTT Failed to connect to broker ", config.Broker)
	} else {
		mqtt.Subscribe("ss/d/+/+", otto.GetDataManager())		
	}

	// start web server / rest server
	server := otto.GetServer()
	go server.Start()

	if interactive {
		go cleanup()
	} else {
		cleanup()
	}
}

func cleanup() {
	<-done
	l.Println("Done, cleaning up()")
}
