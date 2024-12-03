package main

import (
	"fmt"
	"log"

	"github.com/sensorstation/otto"
)

var (
	l *log.Logger
)

func main() {
	l = otto.GetLogger()
	fmt.Printf("logger: %+v\n", l)
	l.Println("Test writting")
	Execute()
}
