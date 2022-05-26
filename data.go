package main

import (
	"log"
	"strings"
	"time"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	consumerQ chan Msg
)

func init() {
	consumerQ = make(chan Msg)
}

// TimeseriesCB call and parse callback data
func dataCB(mc mqtt.Client, mqttmsg mqtt.Message) {
	topic := mqttmsg.Topic()

	// extract the station from the topic
	paths := strings.Split(topic, "/")
	// root	:= paths[0] 
	category:= paths[1] 
	station := paths[2]
	sensor  := paths[3]
	payload := mqttmsg.Payload()

	consumers := hub.GetConsumers(category) 
	if consumers == nil {
		log.Println("DataCB no consumers for ", topic)
		return					// nobody is listening
	}

	log.Printf("MQTT Message topic %s - value %s\n", topic, string(payload))
	switch (category) {
	case "data":
		msg := Msg{}
		msg.Station = station
		msg.Sensor = sensor
		msg.Data = payload
		msg.Time = time.Now().Unix()
		for _, consumer := range consumers {
			consumer.GetRecvQ() <- msg
		}

	default:
		log.Println("Warning: do not know how to handle", topic)
	}
}
