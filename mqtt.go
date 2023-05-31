package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

type Subscriber struct {
	ID string
	Path string
	mqtt.MessageHandler
	Consumers []Consumer
}

func (sub *Subscriber) String() string {
	return sub.ID + " " + sub.Path
}


// Publisher periodically reads from an io.Reader then publishes that value
// to a corresponding channel
type Publisher struct {
	Path   string
	Period time.Duration
	publishing bool
}

func NewPublisher(p string) (pub *Publisher) {
	pub = &Publisher{
		Path:   p,
		Period: 5 * time.Second,
	}
	return pub
}

// Publish will start producing data from the given data producer via
// the q channel returned to the caller. The caller lets Publish know
// to stop sending data when it receives a communication from the done channel
func (p *Publisher) Publish(done chan string) {
	ticker := time.NewTicker(p.Period)

	go func() {
		defer ticker.Stop()
		p.publishing = true
		for p.publishing {
			select {
			case <-done:
				p.publishing = false
				log.Println("Random Data recieved a DONE, returning")
				break

			case <-ticker.C:
				d := "Hello"
				if d != "" {
					if t := mqttc.Publish(p.Path, byte(0), false, d); t == nil {
						if config.Debug {
							log.Printf("%v - I have a NULL token: %s, %+v", mqttc, p.Path, d)
						}
					}
				}
				log.Printf("publish %s -> %+v\n", p.Path, d)
			}
		}
	}()
}


