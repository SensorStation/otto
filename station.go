package iote

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID         string        `json:"id"`
	LastHeard  time.Time     `json:"last-heard"`
	Expiration time.Duration `json:"expiration"` // how long to timeout a station

	ticker *time.Ticker `json:"-"`
	quit   chan bool    `json:"-"`
	mu     sync.Mutex   `json:"-"`
}

// NewStation creates a new Station with an ID as provided
// by the first parameter
func NewStation(id string) (st *Station) {
	st = &Station{
		ID:         id,
		Expiration: 30 * time.Second,
	}
	return st
}

// Update() will append a new data value to the series
// of data points.
func (s *Station) Update(msg *Msg) {
	s.mu.Lock()
	s.LastHeard = msg.Time
	s.mu.Unlock()
}

func (s *Station) Announce() {
	json, err := json.Marshal(s)
	if err != nil {
		log.Printf("ERROR - Station: %s - jsonified %+v", s.ID, err)
		return
	}
	mqtt.Publish("ss/m/station/"+s.ID, string(json))
}

func (s *Station) Advertise(d time.Duration) {
	if s.ticker == nil {
		s.ticker = time.NewTicker(d)
		s.Announce()
	}

	s.quit = make(chan bool)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.Announce()

			case <-s.quit:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop the station from advertising
func (s *Station) Stop() {
	s.quit <- true
}
