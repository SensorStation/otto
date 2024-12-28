package cmd

import (
	"bufio"
	"os"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/logger"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "read otto commands from a file",
	Long:  `Run otto with the commands in the file`,
	Run:   fileRun,
}

func init() {
	rootCmd.AddCommand(fileCmd)
}

func fileRun(cmd *cobra.Command, args []string) {
	l = logger.GetLogger()
	otto.GetConfig().Interactive = true
	fname := args[0]
	file, err := os.Open(fname)
	if err != nil {
		l.Error(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		RunLine(line)
	}
	if err := scanner.Err(); err != nil {
		l.Error(err.Error())
	}
	return
}
