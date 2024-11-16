/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type mqttConfiguration struct {
	Broker string
}

var (
	mqttConfig mqttConfiguration

	// mqttCmd represents the mqtt command
	mqttCmd = &cobra.Command{
		Use:   "mqtt",
		Short: "Configure and interact with MQTT broker",
		Long:  `This command can be used to interact and diagnose an MQTT broker`,
		Run:   RunMQTT,
	}
)

func init() {
	rootCmd.AddCommand(mqttCmd)

	// Here you will define your flags and configuration settings.
	mqttCmd.Flags().StringVarP(&mqttConfig.Broker, "broker", "b", "localhost", "The IP address of the MQTT Broker")
}

func RunMQTT(cmd *cobra.Command, args []string) {

	switch args[0] {
	case "broker":
		switch len(args) {
		case 1:
			fmt.Println(mqttConfig.Broker)

		case 2:
			mqttReconfig(args[1])

		default:
			fmt.Printf("Error too many commands: %p\n", args)
		}
	}

}

func mqttReconfig(broker string) {
	// if mqtt is already connected shut it down

	// if mqtt is not connected then connect to the broker

}
