package main

import (
	"fmt"
	"log"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/cmd"

	"github.com/spf13/cobra"
)

var (
	l *log.Logger
)

var fooCmd = &cobra.Command{
	Use: "foo",
	Run: fooRun,
}

func main() {
	l = otto.GetLogger()
	cmd.Execute()
}

func init() {
	r := cmd.GetRootCmd()
	r.AddCommand(fooCmd)
}

func fooRun(cmd *cobra.Command, args []string) {
	fmt.Println("Rock the foobar")
}
