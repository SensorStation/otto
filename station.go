package iote

import (
	"fmt"
	"log"
	"time"

	"encoding/json"
	"net/http"
)

var (
	Stations StationManager
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID        string    `json:"id"`
	LastHeard time.Time `json:"last-heard"`
	LastMsg   *Msg      `json:"last-Msg"`
}

// NewStation creates a new Station with an ID as provided
// by the first parameter
func NewStation(id string) (st *Station) {
	st = &Station{
		ID: id,
	}
	return st
}

// Update() will append a new data value to the series
// of data points.
func (s *Station) Update(msg *Msg) {
	//	s.LastHeard = msg.Time
	s.LastMsg = msg
}

type StationManager struct {
	Stations map[string]*Station
}

func NewStationManager() (sm *StationManager) {
	sm = &StationManager{}
	sm.Stations = make(map[string]*Station)
	return sm
}

func (sm *StationManager) Get(stid string) *Station {
	st, _ := sm.Stations[stid]
	return st
}

func (sm *StationManager) Add(st string) (station *Station, err error) {
	if sm.Get(st) != nil {
		return nil, fmt.Errorf("Error adding an existing station")
	}
	station = NewStation(st)
	sm.Stations[st] = station
	return station, nil
}

func (sm *StationManager) Update(stid string, data *Msg) {
	var err error
	st := sm.Get(stid)
	if st == nil {
		log.Println("StationManager: Adding new station: ", stid)
		st, err = sm.Add(stid)
		if err != nil {
			log.Println("StationManager: ERROR Adding new station", stid, err)
			return
		}
	}
	st.Update(data)
}

func (sm *StationManager) Count() int {
	return len(sm.Stations)
}

func (sm StationManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stnames := []string{}

	for _, stn := range sm.Stations {
		stnames = append(stnames, stn.ID)
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(stnames)

	case "POST", "PUT":
		http.Error(w, "Not Yet Supported", 401)
	}
}
