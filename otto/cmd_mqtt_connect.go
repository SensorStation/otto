/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/sensorstation/otto"
	"github.com/spf13/cobra"
)

// brokerCmd represents the broker command
var mqttConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the mqtt broker",
	Long:  `Connect to the MQTT broker`,
	Run:   runMQTTConnect,
}

func init() {
	mqttCmd.AddCommand(mqttConnectCmd)
}

func runMQTTConnect(cmd *cobra.Command, args []string) {
	m := otto.GetMQTT()
	if m.Client == nil || !m.IsConnected() {
		err := m.Connect()
		if err != nil {
			fmt.Printf("Failed to connect to mqtt broker: %s: %s\n", mqtt.Broker, err)
		}
	}
}

func mqtt_init() {
	// Start MQTT
	mqtt := otto.GetMQTT()
	mqtt.Connect()

	// ss/d/<station>/<sensor>
	mqtt.Subscribe("ss/d/+/+", otto.GetSensorManager())
}
