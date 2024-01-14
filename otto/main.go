package main

import (
	"flag"
	"log"

	"github.com/sensorstation/otto"
)

var (
	config Configuration
	o      *otto.OttO
)

func main() {
	var e Echo

	flag.Parse()

	o = otto.NewOttO()
	o.Broker = config.Broker
	o.Addr = config.Addr
	o.Appdir = config.Appdir
	o.Plugins = append(o.Plugins, config.Plugin)

	o.Start()
	o.Register("/api/config", config)
	o.Subscribe("ss/c/otto/#", e)
	<-o.Done
}

type Echo struct {
}

func (e Echo) Callback(t string, payload []byte) {
	log.Println(t, string(payload))
}
