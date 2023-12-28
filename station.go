package iote

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID         string        `json:"id"`
	LastHeard  time.Time     `json:"last-heard"`
	Expiration time.Duration `json:"expiration"` // how long to timeout a station

	Sensors  map[string]*Sensor  `json:"sensors"`
	Controls map[string]*Control `json:"controls"`

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
		Sensors:    make(map[string]*Sensor),
		Controls:   make(map[string]*Control),
	}
	return st
}

// Update() will append a new data value to the series
// of data points.
func (s *Station) Update(msg *Msg) {
	s.mu.Lock()
	data := msg.Data.(MsgData)

	if msg.Type == "d" {
		var sensor *Sensor
		var found bool
		if sensor, found = s.Sensors[data.Device]; !found {
			sensor = &Sensor{
				ID: data.Device,
			}
			s.Sensors[data.Device] = sensor
		}
		sensor.Update(msg.Data, msg.Time)
	}
	s.LastHeard = msg.Time
	s.mu.Unlock()
}

func (s *Station) Announce() {
	json, err := json.Marshal(s)
	if err != nil {
		log.Printf("ERROR - Station: %s - jsonified %+v", s.ID, err)
		return
	}
	mqtt.Publish("ss/m/station/"+s.ID, string(json))
}

func (s *Station) Advertise(d time.Duration) {
	if s.ticker == nil {
		s.ticker = time.NewTicker(d)
		s.Announce()
	}

	s.quit = make(chan bool)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.Announce()

			case <-s.quit:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop the station from advertising
func (s *Station) Stop() {
	s.quit <- true
}

func (s Station) String() string {
	return s.ID
}

func (s Station) MarshalJSON() (j []byte, err error) {
	type Sens struct {
		ID    string  `json:"id"`
		Value float64 `json:"value"`
	}

	type Stat struct {
		ID        string    `json:"id"`
		LastHeard time.Time `json:"last-heard"`
		Sensors   []Sens    `json:"sensors"`
	}

	stat := Stat{
		ID:        s.ID,
		LastHeard: s.LastHeard,
	}

	for id, sens := range s.Sensors {
		data := sens.LastValue.(MsgData)

		v, err := strconv.ParseFloat(data.Value.(string), 64)
		if err != nil {
			log.Printf("ERROR StationJSON ParseFloat: %s %+v", data.Value.(string))
			v = -99.99
		}

		sens := Sens{
			ID:    id,
			Value: v,
		}
		stat.Sensors = append(stat.Sensors, sens)
	}

	j, err = json.Marshal(&stat)
	return j, err
}
