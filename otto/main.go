package main

import (
	"flag"
	"fmt"
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
	fmt.Printf("-----> BEFORE: O.Server: %+v\n", o.Server)

	o.Broker = config.Broker
	o.Addr = config.Addr
	o.Appdir = config.Appdir
	o.Plugins = append(o.Plugins, config.Plugin)

	fmt.Printf("-----> AFTER1: O.Server: %+v\n", o.Server)
	o.Start()
	fmt.Printf("-----> AFTER2: O.Server: %+v\n", o.Server)

	o.Register("/api/config", config)
	o.Subscribe("ss/c/otto/#", e)
	<-o.Done
}

type Echo struct {
}

func (e Echo) Callback(t string, payload []byte) {
	log.Println(t, string(payload))
}
