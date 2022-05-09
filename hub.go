package iote
/*

import (
	"log"
	"encoding/json"
	"net/http"
	"periph.io/x/periph"
	"periph.io/x/periph/host"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Hub struct {
	ID   string // MAC address
	Addr string

	*http.Server
	*periph.State
	mqttc *mqtt.Client

	// Stations	map[string]*Station
	Publishers   map[string]*Publisher
	Subscribers  map[string]*Subscriber
	Consumers	 map[string]chan Msg

	Done chan string
}

// NewHub creates a hub that did not previously exist.
// The ID will be populated with the MAC address of this node
func NewHub(cfg *Configuration) (s *Hub) {
	s = &Hub{
		ID:          "0xdeadcafe", // MUST get MAC Address - and WIFI SSID
		Addr:        cfg.Addr,
		Publishers:  make(map[string]*Publisher, 10),
		Subscribers: make(map[string]*Subscriber, 10),
		Done:        make(chan string),
	}

	mqtt_connect()
	if mqttc == nil {
		log.Fatalf("Failed to connect to MQTT broker")
	}

	var err error
	s.State, err = host.Init()
	if err != nil {
		log.Printf("Initializing GPIO failed - no GPIO")
		s.State = nil
	}
	return s
}


// Start the HTTP server and serve up the home web app and
// our REST API
func (s *Hub) Start() error {

	log.Println("Connect to our MQTT broker: ", config.Broker)
	if mqttc == nil {
		log.Fatal("Unable to connect to broker, TODO StandAlone mode")
	}

	log.Println("Subscribers: ", len(s.Subscribers))
	for _, sub := range s.Subscribers {
		log.Println("\t" + sub.String())
	}

	log.Println("Starting publishers: ", len(s.Publishers))
	for _, p := range s.Publishers {
		log.Println("\t" + p.Path)
		p.Publish(s.Done)
	}

	log.Println("Starting hub Web and REST server on ", s.Addr)
	return http.ListenAndServe(s.Addr, nil)
}

func (s *Hub) Subscribe(id string, path string, f mqtt.MessageHandler) {
	sub := &Subscriber{id, path, f, nil}
	s.Subscribers[id] = sub

	qos := 0
	if token := mqttc.Subscribe(path, byte(qos), f); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		if config.Verbose {
			log.Printf("subscribe token: %v", token)
		}
	}
	log.Println(id, " subscribed to ", path)
}

func (s *Hub) AddConsumer(id string, consumer Consumer) {
	sub, ex := s.Subscribers[id]
	if !ex {
		log.Println("Error: AddConsumer - unknown subscriber ", id)
		return
	}
	sub.Consumers = append(sub.Consumers, consumer)
}

func (s *Hub) GetConsumers(id string) []Consumer {
	sub, ex := s.Subscribers[id]
	if !ex {
		log.Println("Unknown subscriber ", id)
		return nil	
	}
	return sub.Consumers
}


// ServeHTTP provides a REST interface to the config structure
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(stations.Stations)
	}
}
*/
