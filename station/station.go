package station

import (
	"sync"
	"time"

	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/messanger"
)

type Controller interface {
	Start() error
	Stop() error
}

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID         string        `json:"id"`
	LastHeard  time.Time     `json:"last-heard"`
	Expiration time.Duration `json:"expiration"` // how long to timeout a station

	devices.DeviceManager `json:"devices"`
	Controller            func() error `json:"-"`

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
func (s *Station) Update(msg *messanger.Msg) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastHeard = time.Now()
}

// Stop the station from advertising
func (s *Station) Stop() {
	s.quit <- true
}

func (s *Station) AddDevice(device any) {
	// s.DeviceManager.Add(device)
}
