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
var mqttSubCmd = &cobra.Command{
	Use:   "sub",
	Short: "Subscribe to the mqtt topic",
	Long:  `Subscribe to mqtt tocpic`,
	Run:   runMQTTSub,
}

func init() {
	mqttCmd.AddCommand(mqttSubCmd)
}

func runMQTTSub(cmd *cobra.Command, args []string) {
	m := otto.GetMQTT()
	if m.Client == nil || !m.IsConnected() {
		err := m.Connect()
		if err != nil {
			fmt.Printf("Failed to connect to mqtt broker: %s: %s\n", m.Broker, err)
		}
	}

	p := &otto.MQTTPrinter{}
	m.Subscribe(args[0], p)
}
