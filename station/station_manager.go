package station

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/sensorstation/otto/messanger"
	"github.com/sensorstation/otto/server"
)

// StationManager keeps track of all the stations we have seen
type StationManager struct {
	Stations map[string]*Station `json:"stations"`
	Stale    map[string]*Station `json:"stale"`
	EventQ   chan *StationEvent  `json:"-"`

	ticker *time.Ticker `json:"-"`
	mu     *sync.Mutex  `json:"-"`
}

type StationEvent struct {
	Type      string `json:"type"`
	Device    string `json:"device"`
	StationID string `json:"stationid"`
	Value     string `json:"value"`
}

var (
	stations *StationManager
)

func GetStationManager() *StationManager {
	if stations == nil {
		stations = NewStationManager()
	}
	return stations
}

func NewStationManager() (sm *StationManager) {
	sm = &StationManager{}
	sm.Stations = make(map[string]*Station)
	sm.Stale = make(map[string]*Station)
	sm.mu = new(sync.Mutex)
	sm.EventQ = make(chan *StationEvent)
	return sm
}

func (sm *StationManager) Callback(msg *messanger.Msg) {
	sm.Update(msg)
}

func (sm *StationManager) Start() {

	srv := server.GetServer()
	srv.Register("/api/stations", sm)

	// Start a ticker to clean up stale entries
	quit := make(chan struct{})
	sm.ticker = time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-sm.ticker.C:
				for id, st := range sm.Stations {

					// Don't try to timeout the local station
					if st.Local {
						continue
					}

					// Do not timeout stations with a duration of 0
					if st.Expiration == 0 {
						slog.Info("Station %s expiration == 0 do not timeout", "id", id)
						continue
					}

					// Timeout a station if we have not heard from it in 3
					// timeframes.
					st.mu.Lock()

					expires := st.LastHeard.Add(st.Expiration)
					if expires.Sub(time.Now()) < 0 {
						sm.mu.Lock()
						slog.Info("Station has timed out", "station", id)
						sm.Stale[id] = st
						delete(sm.Stations, id)
						sm.mu.Unlock()
					}
					st.mu.Unlock()
				}

			case ev := <-sm.EventQ:
				slog.Info("Station Event", "event", ev)
				st := sm.Get(ev.StationID)
				if st == nil {
					slog.Warn("Station Event could not find station", "station", ev.StationID)
					continue
				}

			case <-quit:
				sm.ticker.Stop()
				return
			}
		}
	}()
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

func (sm *StationManager) Update(msg *messanger.Msg) (st *Station) {

	var err error

	stid := msg.Station()
	if stid == "" {
		fmt.Printf("Msg path does not include staionId: %q\n", msg.Path)
		return nil
	}

	st = sm.Get(stid)
	if st == nil {
		if st, err = sm.Add(stid); err != nil {
			fmt.Println("Station Manager failed to create new station: ", stid, err)
			return nil
		}
	}

	// data := msg.Data
	st.Update(msg)
	return st
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
