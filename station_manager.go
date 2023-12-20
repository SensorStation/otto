package iote

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	Stations StationManager
)

func init() {
	Stations = NewStationManager()
}

type StationManager struct {
	Stations map[string]*Station
}

func NewStationManager() (sm StationManager) {
	sm = StationManager{}
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
