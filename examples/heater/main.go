package main

import (
	"flag"
	"log"

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
	}
	mqtt.Start()
	mqtt.Subscribe("meta", "ss/m/#", StationCallback)
	mqtt.Subscribe("data", "ss/d/#", SubscribeCallback)

	srv = iote.Server{
		Addr:   config.Addr,
		Appdir: "/srv/iot/iotvue/dist",
	}
	srv.Register("/api/config", config)
	srv.Start(config.Addr)
}

func StationCallback(mc gomqtt.Client, mqttmsg gomqtt.Message) {
	// log.Printf("Incoming: %s, %q", mqttmsg.Topic(), mqttmsg.Payload())
	msg, err := iote.MsgFromMQTT(mqttmsg.Topic(), mqttmsg.Payload())
	log.Printf("Incoming Station Meta: %+v", msg)

	iote.Stations.Update(msg)
}

// TimeseriesCB call and parse callback msg
func SubscribeCallback(mc gomqtt.Client, mqttmsg gomqtt.Message) {

	// log.Printf("Incoming: %s, %q", mqttmsg.Topic(), mqttmsg.Payload())
	msg := iote.MsgFromMQTT(mqttmsg.Topic(), mqttmsg.Payload())
	log.Printf("Incoming: %+v", msg)

	if msg.Device == "tempc" || msg.Device == "tempf" {
		controller.Update(msg)
	}

	// update the station that sent the msg
	iote.Store.Store(msg)
	iote.Stations.Update(msg.Station, msg)
}

type Controller struct {
	Max float64
	Min float64
}

func (c Controller) On(station string) {
	mqtt.Publish("ss/c/"+station+"header", "on")
}

func (c Controller) Off(station string) {
	mqtt.Publish("ss/c/"+station+"header", "off")
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
