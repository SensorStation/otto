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
	mqttCmd.PersistentFlags().StringVar(&mqttConfig.Broker, "broker", "localhost", "Set the MQTT Broker")
}

func mqttRun(cmd *cobra.Command, args []string) {
	m := otto.GetMQTT()
	if m == nil {
		fmt.Println("MQTT has not been initialized")
		return
	}

	// If the broker config changes and mqtt is connected, disconnect
	// and reconnect to new broker
	if mqttConfig.Broker != m.Broker {
		m.Broker = mqttConfig.Broker
	}

	connected := false
	if m.Client != nil {
		connected = m.IsConnected()
	}

	fmt.Printf("   Broker: %s\n", m.Broker)
	fmt.Printf("Connected: %t\n", connected)
	fmt.Printf("    Debug: %t\n", m.Debug)
	fmt.Printf("Subscribers")

	for _, s := range m.Subscribers {
		fmt.Println(s)
	}
}
