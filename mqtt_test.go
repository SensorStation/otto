package iote

import (
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TestMQTT ensures that we can a) subscribe to a specific channel
// b) publish data to a specific channel and c) recieve the data before
// a timeout.  In this test we 
func TestMQTT(t *testing.T) {
	msg := GetMessanger()
	if msg == nil {
		t.Error("Expected a messanger but got nil")
	}

	topic := "iote/test"
	message := "Hello, World!"
	heard := make(chan bool)
	msg.Subscribe("test", topic, func(c mqtt.Client, m mqtt.Message) {

		// This anonymous function is the callback for all messages sent
		// to the MQTT 'iot/test' topic
		if topic != m.Topic() {
			t.Errorf("Expected topic (%s) got (%s)", topic, m.Topic())
		}
		if message != string(m.Payload()) {
			t.Errorf("Message expected (%s) got (%s) ", message, m.Payload())
		}
		heard <- true		
	})

	msg.Publish("iote/test", message)

	select {
	case <- heard:
		// Our message has been recieved. Yeah the test passed! Say nothing.

	case <-time.After(time.Second * 5):
		t.Error("Expected a message from client got nothing")
	}
}
