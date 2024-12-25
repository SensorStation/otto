package cmd

import (
	"fmt"

	"github.com/sensorstation/otto/devices"
	"github.com/spf13/cobra"
)

var (
	gpioCmd = &cobra.Command{
		Use:   "gpio",
		Short: "Interact with gpio pins",
		Long:  "Configure, read and set GPIO pins on a Raspberry Pi",
		Run:   gpioRun,
	}
	g *devices.GPIO
)

func init() {
	rootCmd.AddCommand(gpioCmd)
}

func gpioRun(cmd *cobra.Command, args []string) {
	g = devices.GetGPIO()
	str := g.String()
	if str == "" {
		str = "GPIO has not been configured"
	}
	fmt.Println(str)
}
