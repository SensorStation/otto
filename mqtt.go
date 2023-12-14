package iote

import (
	"fmt"
	"log"
	"os"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTT struct {
	id     string
	Broker string
	Debug  bool

	subscribers map[string]*Subscriber
	gomqtt.Client
}

func (m *MQTT) Start() {
	m.subscribers = make(map[string]*Subscriber)

	m.Connect()
	m.Subscribe("data", "#", SubscribeCallback)
}

func (m *MQTT) Connect() {
	if m.Debug {
		gomqtt.DEBUG = log.New(os.Stdout, "", 0)
		gomqtt.ERROR = log.New(os.Stdout, "", 0)
	}

	m.id = "sensorStation"
	m.Broker = "tcp://" + m.Broker + ":1883"

	connOpts := gomqtt.NewClientOptions().AddBroker(m.Broker).SetClientID(m.id).SetCleanSession(true)
	m.Client = gomqtt.NewClient(connOpts)
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	log.Println("Connected to broker: ", m.Broker)
}

//
// ss/<ethaddr>/<data>/tempc value
// ss/<ethaddr>/<data>/humidity value
//
func (m *MQTT) Subscribe(id string, path string, f gomqtt.MessageHandler) {
	sub := &Subscriber{id, path, f}
	m.subscribers[id] = sub

	qos := 0
	if token := m.Client.Subscribe(path, byte(qos), f); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		if m.Debug {
			log.Printf("subscribe token: %v", token)
		}
	}
	log.Println(id, "subscribed to", path)
}

// TimeseriesCB call and parse callback msg
func SubscribeCallback(mc gomqtt.Client, mqttmsg gomqtt.Message) {

	// log.Printf("Incoming: %s, %q", mqttmsg.Topic(), mqttmsg.Payload())

	msg := MsgFromMQTT(mqttmsg.Topic(), mqttmsg.Payload())
	log.Printf("Incoming: %+v", msg)
	// disp.InQ <- msg

	// update the station that sent the msg
	// stations.Update(msg.Station, msg)
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
/*
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
*/
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
