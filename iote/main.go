package main

import (
	"flag"
	"fmt"
	
	"github.com/rustyeddy/iote"	
)


var (
	config iote.Configuration
)

func init() {
	config = iote.GetConfig()

	flag.StringVar(&config.Broker, "broker",  "tcp://localhost:1883", "Address of MQTT broker")
}

func main() {
	flag.Parse()

	fmt.Printf("Config: %+v\n", config)

	msg := iote.GetMessanger()
	msg.Subscribe("iote/data", iote.DataCB) // a little weird

	// srv := iote.GetHTTP()
	// srv.Register("/api/config", iote.Config)
}
