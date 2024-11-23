/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// brokerCmd represents the broker command
var brokerCmd = &cobra.Command{
	Use:   "broker",
	Short: "Print or set the MQTT broker",
	Long:  `Print or set the MQTT broker`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("broker called")
	},
}

func init() {
	mqttCmd.AddCommand(brokerCmd)
}

func mqttBroker(cmd *cobra.Command, args []string) {

}
