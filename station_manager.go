package iote

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// StationManager keeps track of all the stations we have seen
type StationManager struct {
	Stations map[string]*Station `json:"stations"`
	Stale    map[string]*Station `json:"stale"`

	ticker *time.Ticker `json:"-"`
	mu     sync.Mutex   `json:"-"`
}

var (
	Stations StationManager
)

func init() {
	Stations = NewStationManager()
}

func NewStationManager() (sm StationManager) {
	sm = StationManager{}
	sm.Stations = make(map[string]*Station)
	sm.Stale = make(map[string]*Station)

	// Start a ticker to clean up stale entries

	quit := make(chan struct{})
	sm.ticker = time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-sm.ticker.C:
				for id, st := range sm.Stations {

					// Do not timeout stations with a duration of 0
					if st.Expiration == 0 {
						log.Printf("Station %s expiration == 0 do not timeout", id)
						continue
					}

					// Timeout a station if we have not heard from it in 3
					// timeframes.
					st.mu.Lock()

					expires := st.LastHeard.Add(st.Expiration)
					if expires.Sub(time.Now()) < 0 {
						sm.mu.Lock()
						log.Printf("Station: %s has timed out\n", id)
						sm.Stale[id] = st
						delete(sm.Stations, id)
						sm.mu.Unlock()
					}
					st.mu.Unlock()
				}

			case <-quit:
				sm.ticker.Stop()
				return
			}
		}
	}()

	return sm
}

func (sm *StationManager) Get(stid string) *Station {
	sm.mu.Lock()
	st, _ := sm.Stations[stid]
	sm.mu.Unlock()
	return st
}

func (sm *StationManager) Add(st string) (station *Station, err error) {
	if sm.Get(st) != nil {
		return nil, fmt.Errorf("Error adding an existing station")
	}
	station = NewStation(st)
	sm.mu.Lock()
	sm.Stations[st] = station
	sm.mu.Unlock()
	return station, nil
}

func (sm *StationManager) Update(msg *Msg) {
	var err error
	st := sm.Get(msg.Station)
	if st == nil {
		log.Println("StationManager: Adding new station: ", msg.Station)
		st, err = sm.Add(msg.Station)
		if err != nil {
			log.Println("StationManager: ERROR Adding new station", msg.Station, err)
			return
		}
	}
	st.Update(msg)
}

func (sm *StationManager) Count() int {
	return len(sm.Stations)
}

func (sm StationManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(sm)

	case "POST", "PUT":
		http.Error(w, "Not Yet Supported", 401)
	}
}
