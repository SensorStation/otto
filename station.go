package otto

import (
	"sync"
	"time"
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID         string        `json:"id"`
	LastHeard  time.Time     `json:"last-heard"`
	Expiration time.Duration `json:"expiration"` // how long to timeout a station

	*DataManager 

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
	st.DataManager = NewDataManager()
	return st
}

// Update() will append a new data value to the series
// of data points.
func (s *Station) Update(msg *Msg) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastHeard = time.Now()
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

	if s.DataManager.DataMap[label] == nil {
		s.DataManager.DataMap[label] = NewTimeseries(label)
	}

	s.DataManager.DataMap[label].Data = append(s.DataMap[label].Data, d)
}
