package main

import (
	"log"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/cmd"
)

var (
	l *log.Logger
)

func main() {
	l = otto.GetLogger()
	cmd.Execute()
}
