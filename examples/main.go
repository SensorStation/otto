package main

import (
	"flag"
	"net/http"
	"sync"

	"github.com/rustyeddy/iote"
	// "net/http"
)

// Globals
var (
	config Configuration
	mqtt   iote.MQTT

	disp  *iote.Dispatcher
	srv   *iote.Server
	wserv iote.Websock
)

func init() {
	disp = iote.NewDispatcher()
}

func main() {
	var wg sync.WaitGroup

	// Parse command line argumens and update the config as appropriate
	flag.Parse()

	mqtt = iote.MQTT{
		Broker: config.Broker,
	}
	mqtt.Start()

	// The web app
	fs := http.FileServer(http.Dir("/srv/iot/iotvue/dist"))
	// Now create the station based on the given configuration
	srv = iote.NewServer(config.Addr)
	srv.Register("/", fs)
	srv.Register("/ws", wserv)
	srv.Register("/ping", Ping{})
	srv.Register("/api/config", config)
	srv.Register("/api/data", srv)
	srv.Register("/api/stations", iote.Stations)

	wg.Add(1)
	go srv.Start(config.Addr, wg)

	wg.Wait()
}
