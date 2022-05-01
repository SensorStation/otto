package iote

import (
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (

)

func TestMQTT(t *testing.T) {
	msg := GetMessanger()
	if msg == nil {
		t.Error("Expected a messanger but got nil")
	}

	topic := "iote/test"
	message := "Hello, World!"
	heard := make(chan bool)
	msg.Subscribe("test", topic, func(c mqtt.Client, m mqtt.Message) {
		if topic != m.Topic() {
			t.Errorf("Expected topic (%s) got (%s)", topic, m.Topic())
		}
		if message != string(m.Payload()) {
			t.Errorf("Message expected (%s) got (%s) ", message, m.Payload())
		}
		heard <- true		
	})

	msg.Publish("iote/test", 0, false, message)
	select {
	case <- heard:
		// Our message has been recieved. 
	case <-time.After(time.Second * 5):
		t.Error("Expected a message from client got nothing")
	}
}
