package main

import (
	"fmt"
	"log"
	"os"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttc mqtt.Client
)

func mqtt_connect() {
	if config.DebugMQTT {
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
	}

	id := "sensorStation"
	broker := "tcp://" + config.Broker + ":1883"
	
	connOpts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(id).SetCleanSession(true)
	mqttc = mqtt.NewClient(connOpts)
	if token := mqttc.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

type ToCloudMsg struct {
	Addr	string `json:"addr"`
	Type	string `json:"type"`
	Data	map[string]interface{}	`json:"data"`
}

// ToCloudCB is the callback when we recieve MQTT messages on the '/mesh/xxxxxx/toCloud' channel. 
func ToCloudCB(mc mqtt.Client, msg mqtt.Message) {
	if false {
		log.Printf("Incoming message topic: %s\n", msg.Topic());
	}
	mesh.MsgRecv(msg.Topic(), msg.Payload())
}



