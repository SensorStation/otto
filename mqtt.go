package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	ID     string
	Broker string

	Publishers  map[string]*Publisher
	Subscribers map[string]*Subscriber

	gomqtt.Client
}

func NewMQTT() *MQTT {
	return &MQTT{
		ID:          "IoTe",
		Broker:      config.Broker,
		Publishers:  make(map[string]*Publisher),
		Subscribers: make(map[string]*Subscriber),
	}
}

func (m *MQTT) Connect() {
	if config.DebugMQTT {
		gomqtt.DEBUG = log.New(os.Stdout, "", 0)
		gomqtt.ERROR = log.New(os.Stdout, "", 0)
	}

	m.ID = "sensorStation"
	m.Broker = "tcp://" + config.Broker + ":1883"

	connOpts := gomqtt.NewClientOptions().AddBroker(m.Broker).SetClientID(m.ID).SetCleanSession(true)
	m.Client = gomqtt.NewClient(connOpts)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	log.Println("Connected to broker: ", m.Broker)
}

func (m *MQTT) Subscribe(id string, path string, f gomqtt.MessageHandler) {
	sub := &Subscriber{id, path, f}
	m.Subscribers[id] = sub

	qos := 0
	if token := m.Client.Subscribe(path, byte(qos), f); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		if config.Verbose {
			log.Printf("subscribe token: %v", token)
		}
	}
	log.Println(id, "subscribed to", path)
}

// TimeseriesCB call and parse callback msg
func msgCB(mc gomqtt.Client, mqttmsg gomqtt.Message) {
	topic := mqttmsg.Topic()

	// extract the station from the topic
	paths := strings.Split(topic, "/")

	// ss/msg/<source>/<sensor> <value>
	msg := &Msg{
		Source:   paths[1],
		Category: paths[2],
		Device:   paths[3],
		Time:     time.Now(),
	}

	var val float64
	fmt.Sscanf(string(mqttmsg.Payload()), "%f", &val)
	msg.Value = val

	disp.InQ <- msg

	// update the station that sent the msg
	stations.Update(msg.Source, msg)
}

type Subscriber struct {
	ID   string
	Path string
	gomqtt.MessageHandler
}

func (sub *Subscriber) String() string {
	return sub.ID + " " + sub.Path
}

// Publisher periodically reads from an io.Reader then publishes that value
// to a corresponding channel
type Publisher struct {
	Path       string
	Period     time.Duration
	publishing bool
}

func NewPublisher(p string) (pub *Publisher) {
	pub = &Publisher{
		Path:   p,
		Period: 5 * time.Second,
	}
	return pub
}

// Publish will start producing msg from the given data producer via
// the q channel returned to the caller. The caller lets Publish know
// to stop sending data when it receives a communication from the done channel
func (m MQTT) Publish(done chan string) {
	// ticker := time.NewTicker(p.Period)

	// go func() {
	// 	defer ticker.Stop()
	// 	p.publishing = true
	// 	for p.publishing {
	// 		select {
	// 		case <-done:
	// 			p.publishing = false
	// 			log.Println("Random Data recieved a DONE, returning")
	// 			break

	// 		case <-ticker.C:
	// 			d := "Hello"
	// 			if d != "" {
	// 				if t := m.Client.Publish(p.Path, byte(0), false, d); t == nil {
	// 					if config.Debug {
	// 						log.Printf("%v - I have a NULL token: %+v", m.Client, p.Path, d)
	// 					}
	// 				}
	// 			}
	// 			log.Printf("publish %s -> %+v\n", p.Path, d)
	// 		}
	// 	}
	// }()
}
