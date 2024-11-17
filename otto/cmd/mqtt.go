/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
		Run:   RunMQTT,
	}
)

func init() {
	rootCmd.AddCommand(mqttCmd)
}

func RunMQTT(cmd *cobra.Command, args []string) {

	fmt.Println("TODO print mqtt configuration and connectivity")
}

func mqttDisable() {

}

func mqttReconfig(broker string) {
	// if mqtt is already connected shut it down

	// if mqtt is not connected then connect to the broker

}
