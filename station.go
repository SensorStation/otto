package otto

import (
	"fmt"
	"sync"
	"time"
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID         string        `json:"id"`
	LastHeard  time.Time     `json:"last-heard"`
	Expiration time.Duration `json:"expiration"` // how long to timeout a station

	Timeseries map[string]*Timeseries

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
	st.Timeseries = make(map[string]*Timeseries)
	return st
}

// Update() will append a new data value to the series
// of data points.
func (s *Station) Update(msg *Msg) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, err := time.Parse(time.RFC3339, msg.Time.String())
	if err != nil {
		l.Println("Station Failed to parse msg.Time:", err)
	} else {
		s.LastHeard = t
	}
}

// Stop the station from advertising
func (s *Station) Stop() {
	s.quit <- true
}

func (s *Station) Insert(label string, val interface{}) {
	d := &Data{
		Time:  time.Now(),
		Value: val,
	}

	fmt.Printf("TS: %+v\n", s)

	if s.Timeseries[label] == nil {
		s.Timeseries[label] = NewTimeseries(label)
	}

	s.Timeseries[label].Data = append(s.Timeseries[label].Data, d)
}
