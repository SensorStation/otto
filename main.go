package main

import (
	"flag"

	"net/http"
)

// Globals
var (
	hub		 *Hub
	mesh	 MeshNetwork
	stations StationManager
	wserv	 WSServer
)

func init() {
	mesh = MeshNetwork{
		Nodes: make(map[string]*MeshNode),
	}
	stations = NewStationManager()
}

func main() {

	// Parse command line argumens and update the config as appropriate
	flag.Parse()

	// Create the state configuration for this station.
	cfg := GetConfig()

	// Now create the station based on the given configuration
	hub = NewHub(&cfg)

	// Register our REST callbacks, specifically answer to pings
	hub.Register("/ws", wserv)
	hub.Register("/ping", Ping{})
	hub.Register("/api/config", config)
	hub.Register("/api/mesh", mesh)
	hub.Register("/api/data", hub)
	hub.Register("/api/stations", stations)

	// The web app
	fs := http.FileServer(http.Dir("/srv/iot/iotvue/dist"))
	hub.Register("/", fs)

	// Subscribe to MQTT channels
	hub.Subscribe("mesh", "mesh/+/toCloud", ToCloudCB)
	hub.Subscribe("net",  "ss/net/announce", ToCloudCB)
	hub.Subscribe("data", "ss/data/+/+", dataCB)

	// Add the Stations Consumer for in memory copies
	hub.AddConsumer("data", stations)
	go stations.Listen()

	// ----------------------------------------------------------
	// Register our publishers with their respective readers
	// ----------------------------------------------------------
	if config.Mock {
		//	pub := NewPublisher("data/cafedead/tempf", hub.NewRando())		
		//	AddPublisher("data/cafedead/humidity", hub.NewRando())		
	}
	hub.Start()
}

