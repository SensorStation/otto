/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/sensorstation/otto"
	"github.com/spf13/cobra"
)

// brokerCmd represents the broker command
var mqttPubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Publish to the mqtt topic",
	Long:  `Publish to mqtt tocpic`,
	Run:   runMQTTPub,
}

func init() {
	mqttCmd.AddCommand(mqttPubCmd)
}

func runMQTTPub(cmd *cobra.Command, args []string) {
	m, err := otto.GetMQTT()
	if err != nil {
		fmt.Println("MQTT Is not connected to a broker")
		return
	}

	m.Publish(args[0], args[1])
}
