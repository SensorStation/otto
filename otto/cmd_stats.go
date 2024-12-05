package main

import (
	"fmt"

	"github.com/sensorstation/otto"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use: "stats",
	Short: "Display runtime stats",
	Long: `Display runtime stats`,
	Run: statsRun,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func statsRun(cmd *cobra.Command, args []string) {
	stats := otto.GetStats()
	fmt.Printf("Stats: %+v\n", stats)
}
