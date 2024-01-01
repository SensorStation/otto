package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rustyeddy/iote"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

// Globals
var (
	config     Configuration
	controller Controller

	disp *iote.Dispatcher
	mqtt iote.MQTT
	srv  iote.Server
)

func main() {

	// Parse command line argumens and update the config as appropriate
	flag.Parse()

	disp = iote.GetDispatcher()
	fmt.Printf("Dispatcher: %+v\n", disp)

	mqtt = iote.MQTT{
		Broker: config.Broker,
		ID:     "heater",
	}
	mqtt.Start()
	time.Sleep(1 * time.Second)
	mqtt.Subscribe("data", "ss/d/#", DataCallback)

	srv = iote.Server{
		Addr:   config.Addr,
		Appdir: "/srv/iot/iotvue/dist",
	}

	err := srv.Start(config.Addr)
	if err != nil {
		log.Printf("ERROR HTTP Server: %+v", err)
		os.Exit(1)
	}
	log.Printf("%s shutting down", os.Args[0])
}

// TimeseriesCB call and parse callback msg
func DataCallback(mc gomqtt.Client, mqttmsg gomqtt.Message) {

	log.Printf("MQTT [I] Data: %s, %s", mqttmsg.Topic(), mqttmsg.Payload())
	msg, err := iote.MsgFromMQTT(mqttmsg.Topic(), mqttmsg.Payload())
	if err != nil {
		log.Printf("ERROR - parsing incoming message: %+v\n", err)
		return
	}

	controller.Update(&msg.Data)

	// update the station that sent the msg
	// iote.Store.Store(msg)
	station := iote.Stations.Update(msg)
	if station == nil {
		log.Printf("Failed to update station for %+v\n", msg)
		return
	}

	msg.Type = "station"
	disp.InQ <- msg
}
