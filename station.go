package iote

import (
	"encoding/json"
	"log"
	"time"
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID            string    `json:"id"`
	LastHeard     time.Time `json:"last-heard"`
	time.Duration `json:"duration"`

	ticker *time.Ticker `json:"-"`
}

// NewStation creates a new Station with an ID as provided
// by the first parameter
func NewStation(id string) (st *Station) {
	st = &Station{
		ID:       id,
		Duration: 5,
	}
	return st
}

// Update() will append a new data value to the series
// of data points.
func (s *Station) Update(msg *Msg) {
	log.Println("Updating station: ", s.ID)
	s.LastHeard = msg.Time
}

func (s *Station) Announce() {
	json, err := json.Marshal(s)
	if err != nil {
		log.Printf("ERROR - Station: %s - jsonified %+v", s.ID, err)
		return
	}
	mqtt.Publish("ss/m/"+s.ID+"/station", string(json))
}

func (s *Station) Advertise(d time.Duration) {
	s.Duration = d
	if s.ticker == nil {
		s.ticker = time.NewTicker(d * time.Second)
	}

	s.Announce()

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.Announce()

			case <-quit:
				s.ticker.Stop()
				return
			}
		}
	}()
}
