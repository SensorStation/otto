package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "otto",
	Short: "OttO is an IoT platform for creating cool IoT apps and hubs",
	Long: `This is cool stuff and you will be able to find a lot of cool information 
                in the following documentation https://rustyeddy.com/otto/`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Printf("Start otto here: %q\n", args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
