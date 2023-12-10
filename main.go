package main

import (
	"flag"
	"net/http"
	"sync"
	// "net/http"
)

// Globals
var (
	config   Configuration
	disp     *dispatcher
	mqtt     *MQTT
	stations *StationManager
	srv      *Server
	wserv    websock
)

func init() {
	disp = newDispatcher()
}

func main() {
	var wg sync.WaitGroup

	// Parse command line argumens and update the config as appropriate
	flag.Parse()

	// Subscribe to MQTT channels
	// hub = NewHub(&cfg)
	mqtt = NewMQTT()
	mqtt.Connect()
	// mqtt.Subscribe("data", "ss/+/+/+", msgCB)
	mqtt.Subscribe("data", "#", msgCB)

	// Add the Stations Consumer for in memory copies
	// hub.AddConsumer("data", stations)
	stations = NewStationManager()

	// The web app
	fs := http.FileServer(http.Dir("/srv/iot/iotvue/dist"))
	// Now create the station based on the given configuration
	srv = NewServer(config.Addr)
	srv.Register("/", fs)
	srv.Register("/ws", wserv)
	srv.Register("/ping", Ping{})
	srv.Register("/api/config", config)
	srv.Register("/api/data", srv)
	srv.Register("/api/stations", stations)

	wg.Add(1)
	go srv.Start(config.Addr, wg)

	wg.Wait()
}
