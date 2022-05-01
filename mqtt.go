package iote

import (
	"fmt"
	"log"
	"os"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Messanger allows a message consumer to provide a callback
// to specific message paths
type Messanger struct {
	mqtt.Client
}

var (
	msg *Messanger 
)

func GetMessanger() *Messanger {
	if (msg == nil) {
		msg = &Messanger{}
		msg.Connect()
	}
	return msg
}

func (m *Messanger) Connect() {
	// config.DebugMQTT = true
	if config.DebugMQTT {
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
	}

	id := "IoTe"
	connOpts := mqtt.NewClientOptions().AddBroker(config.Broker).SetClientID(id).SetCleanSession(true)

	m.Client = mqtt.NewClient(connOpts)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

func (m *Messanger) Subscribe(id string, path string, f mqtt.MessageHandler) {
	// sub := &Subscriber{id, path, f, nil}
	// s.Subscribers[id] = sub

	qos := 0
	if token := m.Client.Subscribe(path, byte(qos), f); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		if config.Verbose {
			log.Printf("subscribe token: %v", token)
		}
	}
}

