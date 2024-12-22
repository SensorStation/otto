/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	otto.GetMQTT()
}
