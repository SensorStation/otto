package station

import (
	"sync"
	"time"

	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/messanger"
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID         string        `json:"id"`
	LastHeard  time.Time     `json:"last-heard"`
	Expiration time.Duration `json:"expiration"` // how long to timeout a station
	IPAddr     string        `json:"ipaddr"`
	MACAddr    string        `json:"macaddr"`

	*messanger.Messanger

	device.DeviceManager `json:"devices"`

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
		Messanger:  messanger.NewMessanger(id),
	}
	return st
}

// Start the station timeout timer or advertisement timer
func (s *Station) Start() {
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

func (s *Station) AddDevice(device device.Name) {
	s.DeviceManager.Add(device)
}

func (s *Station) GetDevice(name string) any {
	d, _ := s.DeviceManager.Get(name)
	return d
}
