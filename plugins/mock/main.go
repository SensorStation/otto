package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/sensorstation/otto"
)

type Station struct {
	ID           string `json:"id"`
	*time.Ticker `json:"-"`
	quit         chan bool
}

var (
	mqtt *otto.MQTT
)

func main() {
	mqtt = &iote.MQTT{
		Broker: "localhost",
		ID:     "mock",
	}
	mqtt.Start()

	stations := make(map[string]*Station)
	ids := []string{
		"sta1",
		"sta2",
		"sta3",
		"sta4",
		"sta5",
	}
	for _, id := range ids {
		st := &Station{
			ID:   id,
			quit: make(chan bool),
		}
		st.Ticker = time.NewTicker(5 * time.Second)
		stations[id] = st
		st.Announce()
		go func() {
			for {
				select {
				case <-st.Ticker.C:
					st.Announce()

				case <-st.quit:
					st.Ticker.Stop()

				}

			}
		}()

	}

	// Get stations and check the count should be 5

	// turn off station 5 to verify it gets timed out
	log.Println("Stopping station 5")
	if st5, found := stations["sta5"]; found {
		st5.Stop()
	} else {
		log.Printf("ERROR Could not get station sta5 to stop it")
	}

	// Get stations again and check the count should be 4, station 5
	// should be missing

	doneQ := make(chan bool)
	<-doneQ
}

func (st *Station) Announce() {
	json, err := json.Marshal(st)
	if err != nil {
		log.Printf("ERROR - Station: %s - jsonified %+v", st.ID, err)
		return
	}
	mqtt.Publish("ss/m/station/"+st.ID, string(json))
}
