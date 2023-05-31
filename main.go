package main

import (
	"flag"
	"sync"
	"net/http"
)

// Globals
var (
	hub		 *Hub
	mesh	 MeshNetwork
	stations StationManager
	wserv	 WSServer
	srv		 *Server
)

func init() {
	mesh = MeshNetwork{
		Nodes: make(map[string]*MeshNode),
	}
	stations = NewStationManager()
}

func main() {
	var wg sync.WaitGroup

	// Parse command line argumens and update the config as appropriate
	flag.Parse()

	// Create the state configuration for this station.
	cfg := GetConfig()

	// Now create the station based on the given configuration
	srv = NewServer(cfg.Addr)

	// Register our REST callbacks, specifically answer to pings
	srv.Register("/ws", wserv)
	srv.Register("/ping", Ping{})
	srv.Register("/api/config", config)
	srv.Register("/api/mesh", mesh)
	srv.Register("/api/data", srv)
	srv.Register("/api/stations", stations)

	// The web app
	fs := http.FileServer(http.Dir("/srv/iot/iotvue/dist"))
	srv.Register("/", fs)

	wg.Add(1)
	go srv.Start(cfg.Addr, wg)

	// Subscribe to MQTT channels
	hub = NewHub(&cfg)
	// hub.Subscribe("mesh", "mesh/+/toCloud", ToCloudCB)
	// hub.Subscribe("net",  "iote/net/announce", ToCloudCB)
	hub.Subscribe("data", "ss/data/+/+", dataCB)

	// Add the Stations Consumer for in memory copies
	hub.AddConsumer("data", stations)

	wg.Add(1)
	go stations.Listen(wg)

	// ----------------------------------------------------------
	// Register our publishers with their respective readers
	// ----------------------------------------------------------
	if config.Mock {
		//	pub := NewPublisher("data/cafedead/tempf", hub.NewRando())		
		//	AddPublisher("data/cafedead/humidity", hub.NewRando())		
	}
	wg.Add(1)
	hub.Start(wg)

	wg.Wait()
}

