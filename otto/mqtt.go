/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/sensorstation/otto"
	"github.com/spf13/cobra"
)

type mqttConfiguration struct {
	Broker  string
	Enabled bool
}

var (
	mqttConfig mqttConfiguration

	// mqttCmd represents the mqtt command
	mqttCmd = &cobra.Command{
		Use:   "mqtt",
		Short: "Configure and interact with MQTT broker",
		Long:  `This command can be used to interact and diagnose an MQTT broker`,
		Run:   mqttRun,
	}
)

func init() {
	rootCmd.AddCommand(mqttCmd)
}

func mqttRun(cmd *cobra.Command, args []string) {
	fmt.Println("TODO print mqtt configuration and connectivity")
	m := otto.GetMQTT()
	if m == nil {
		fmt.Println("MQTT is nil")
	}
	fmt.Printf("MQTT: %+v\n", m)
}
