package main

import (
	"flag"

	"github.com/SensorStation/otto"
)

var (
	config Configuration
)

func main() {
	flag.Parse()

	otto.Init()
}
