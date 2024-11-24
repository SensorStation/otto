/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/sensorstation/otto"
	"github.com/spf13/cobra"
)

var (
	count int

	// mqttCmd represents the mqtt command
	mqttBlastCmd = &cobra.Command{
		Use:   "blast",
		Short: "Start blasting MQTT messages from (count) blasters",
		Long:  `Start blasting MQTT messages from (count) blasters`,
		Run:   mqttBlastRun,
	}
)

func init() {
	mqttCmd.AddCommand(mqttBlastCmd)
	mqttBlastCmd.PersistentFlags().IntVar(&count, "count", 10, "The number of blasters to start")
}

func mqttBlastRun(cmd *cobra.Command, args []string) {
	blasters := otto.GetMQTTBlasters(10)
	blasters.Blast()
}
