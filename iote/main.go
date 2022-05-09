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

	srv := iote.Server
	srv.Register("/api/config", config)
	err := srv.Listen()

	fmt.Printf("Good bye: %+v\n", err)
}
