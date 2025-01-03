package cmd

import (
	"github.com/sensorstation/otto/logger"
	"github.com/spf13/cobra"
)

var (
	l      *logger.Logger
	appdir string
)

var rootCmd = &cobra.Command{
	Use:   "otto",
	Short: "OttO is an IoT platform for creating cool IoT apps and hubs",
	Long: `This is cool stuff and you will be able to find a lot of cool information 
                in the following documentation https://rustyeddy.com/otto/`,
	Run: ottoRun,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&appdir, "appdir", "embed", "root of the web app")
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func Execute() {
	l = logger.GetLogger()
	if err := rootCmd.Execute(); err != nil {
		l.Error(err.Error())
		return
	}
}

func ottoRun(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
