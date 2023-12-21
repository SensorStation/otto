package main

import (
	"github.com/rustyeddy/iote"
)

func main() {
	mqtt := iote.MQTT{
		Broker: "localhost",
		ID:     "mock",
	}
	mqtt.Start()

	stations := make(map[string]*iote.Station)
	ids := []string{
		"sta1",
		"sta2",
		"sta3",
		"sta4",
		"sta5",
	}
	for _, id := range ids {
		st := iote.NewStation(id)
		st.Advertise(5)
		stations[id] = st
	}

	doneQ := make(chan bool)
	<-doneQ

}
