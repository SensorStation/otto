package main

import (
	"flag"
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

	mqtt iote.MQTT
	srv  iote.Server
)

func main() {

	// Parse command line argumens and update the config as appropriate
	flag.Parse()

	mqtt = iote.MQTT{
		Broker: config.Broker,
		ID:     "heater",
	}
	mqtt.Start()
	time.Sleep(1 * time.Second)
	mqtt.Subscribe("meta", "ss/m/station/#", StationCallback)
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

func StationCallback(mc gomqtt.Client, mqttmsg gomqtt.Message) {
	log.Printf("Incoming Station: %s, %s", mqttmsg.Topic(), mqttmsg.Payload())
	msg, err := iote.MsgFromMQTT(mqttmsg.Topic(), mqttmsg.Payload())
	if err != nil {
		log.Printf("ERROR - parsing incoming message: %+v\n", err)
		return
	}
	iote.Stations.Update(msg)
}

// TimeseriesCB call and parse callback msg
func DataCallback(mc gomqtt.Client, mqttmsg gomqtt.Message) {

	log.Printf("Incoming Data: %s, %s", mqttmsg.Topic(), mqttmsg.Payload())
	msg, err := iote.MsgFromMQTT(mqttmsg.Topic(), mqttmsg.Payload())
	if err != nil {
		log.Printf("ERROR - parsing incoming message: %+v\n", err)
		return
	}

	if msg.Device == "tempc" || msg.Device == "tempf" {
		controller.Update(msg)
	}

	// update the station that sent the msg
	iote.Store.Store(msg)
	iote.Stations.Update(msg)
}

type Controller struct {
	Max float64
	Min float64
}

func (c Controller) On(station string) {
	mqtt.Publish("ss/c/"+station+"/heater", "on")
}

func (c Controller) Off(station string) {
	mqtt.Publish("ss/c/"+station+"/heater", "off")
}

func (c Controller) Update(msg *iote.Msg) {

	var v float64

	switch msg.Value.(type) {
	case float64:
		v = msg.Value.(float64)
	}

	if v <= c.Min {
		c.On(msg.Device)
	} else if v >= c.Max {
		c.Off(msg.Device)
	}
}
