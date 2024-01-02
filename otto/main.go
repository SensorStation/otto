package main

import (
	"flag"
	"log"
	"os"

	"github.com/sensorstation/otto"
)

var (
	config Configuration
)

func main() {
	var e Echo

	flag.Parse()
	o := otto.O()

	o.Config = config
	o.Start(config, os.Args[1:])
	o.Register("/api/config", config)
	o.Subscribe("ss/c/otto/#", e)
	<-o.Done
}

type Echo struct {
}

func (e Echo) Callback(t string, payload []byte) {
	log.Println(t, string(payload))
}
