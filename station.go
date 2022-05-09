package iote

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"net/http"
	"encoding/json"
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID			string					`json:"id"`
	LastTime	time.Time			`json:"last-time"`
	Sensors map[string]*Timeseries	`json:"sensors"`


}

var (
	stations map[string]*Station
)

// NewStation creates a new Station with an ID as provided
// by the first parameter
func NewStation(id string) (st *Station) {
	st = &Station{
		ID: id,
	}
	st.Sensors = make(map[string]*Timeseries)
	return st
}

// GetSensor returns the sensor with the matching ID.
// Nil "" will be returned if no sensor with that ID is
// found
func (s *Station) GetSensor(sensor string) *Timeseries {
	var sens *Timeseries
	var e  bool
	
	if sens, e = s.Sensors[sensor]; !e {
		return nil
	}
	return sens
}

// AddData will append another timestamped data to a sensors
// timeseries database.
func (s *Station) AddData(sensor string, val float64) {
	ts := s.GetSensor(sensor)
	if ts == nil {
		ts = NewTimeseries()
		s.Sensors[sensor] = ts
	}
	ts.Append(val)
}

// Length() return the number of sensors handled by this package 
func (s *Station) Length() int {
	return len(s.Sensors)
}

// DataCound() returns the number of Data items held by the
// identified sensor
func (s *Station) DataCount(sensor string) int {
	ts := s.GetSensor(sensor)
	return len(ts.Values)
}

// Update() will append a new data value to the series
// of data points.
func (s *Station) Update(sensor string, data []byte) {

	val, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
	 	log.Println("ERROR Station.Update - ", s.ID, sensor, err)
	 	return 
	}

	sens, e := s.Sensors[sensor]
	if !e {
		log.Println("Adding new Timeseries sensor", sensor)
		sens = NewTimeseries()
		s.Sensors[sensor] = sens
	}
	ts := NewTimestamp(val)
	sens.Values = append(sens.Values, ts)

	if (len(sens.Values) % config.MaxData == 0) {
		_, sens.Values = sens.Values[0], sens.Values[1:]
	}
}

type StationManager struct {
	Stations map[string]*Station
	RecvQ	 chan Msg
}

func NewStationManager() (sm StationManager) {
	sm = StationManager{}
	sm.Stations = make(map[string]*Station)
	sm.RecvQ = make(chan Msg)
	return sm
}

func (sm StationManager) GetID() string {
	return "Station Manager"
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
	sm.Stations[st] = station;
	return station, nil
}

func (sm *StationManager) Update(stid string, sensor string, data []byte) {
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
	st.Update(sensor, data)
}

func (sm *StationManager) Count() int {
	return len(sm.Stations)
}

func (sm StationManager) Recv(msg Msg) {
	sm.Update(msg.Station, msg.Sensor, msg.Data)
}

func (sm StationManager) GetRecvQ() chan Msg {
	return sm.RecvQ
}

func (sm StationManager) Listen() {
	log.Printf("Station Manager Listening on Q %+v", sm.RecvQ)
	for true {
		select {
		case msg := <- sm.RecvQ:
			sm.Recv(msg)
		}
	}
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
