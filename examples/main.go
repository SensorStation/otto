package main

import (
	"flag"
	"sync"

	"github.com/rustyeddy/iote"
)

// Globals
var (
	config Configuration
	mqtt   iote.MQTT
	srv    iote.Server
)

func main() {
	var wg sync.WaitGroup

	// Parse command line argumens and update the config as appropriate
	flag.Parse()

	mqtt = iote.MQTT{
		Broker: config.Broker,
	}
	mqtt.Start()

	srv = iote.Server{
		Addr:   config.Addr,
		Appdir: "/srv/iot/iotvue/dist",
	}
	srv.Register("/api/config", config)
	wg.Add(1)
	srv.Start(config.Addr, wg)

	/*
		go srv.Start(config.Addr, wg)
			wg.Wait()
	*/
}
