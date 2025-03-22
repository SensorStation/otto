package station

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/messanger"
	"github.com/sensorstation/otto/server"
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
// by the first parameter. Here we need to detect a duplicate
// station before trying to register another one.
func NewStation(id string) (st *Station) {
	st = &Station{
		ID:         id,
		Expiration: 30 * time.Second,
		Messanger:  messanger.NewMessanger(id),
	}
	return st
}

func (st *Station) Register() {
	// this needs to move
	srv := server.GetServer()
	srv.Register("/api/station/"+st.ID, st)
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

// AddDevice will do what it says by placing the device with a given
// name in the stations device manager. This library is basically a
// key value store, anything supporting the Name Interface:
// i.e. Name() string.
func (s *Station) AddDevice(device device.Name) {
	s.DeviceManager.Add(device)
}

// GetDevice returns the device (anythig supporting the Name (Name()) interface)
func (s *Station) GetDevice(name string) any {
	d, _ := s.DeviceManager.Get(name)
	return d
}

// Create an endpoint for this device to be queried.
func (s Station) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(s)

	case "POST", "PUT":
		http.Error(w, "Not Yet Supported", 401)
	}
}
