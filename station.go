package otto

import (
	"log"
	"sync"
	"time"
)

// Station is the primary structure that holds an array of
// Sensors which in turn hold a timeseries of datapoints.
type Station struct {
	ID         string        `json:"id"`
	LastHeard  time.Time     `json:"last-heard"`
	Expiration time.Duration `json:"expiration"` // how long to timeout a station

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
func (s *Station) Update(msg *Msg) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, err := time.Parse(time.RFC3339, msg.Time.String())
	if err != nil {
		log.Println("Station Failed to parse msg.Time , err")
	} else {
		s.LastHeard = t
	}
}

// Stop the station from advertising
func (s *Station) Stop() {
	s.quit <- true
}

// func (s *Station) Announce() {
// 	json, err := json.Marshal(s)
// 	if err != nil {
// 		log.Printf("ERROR - Station: %s - jsonified %+v", s.ID, err)
// 		return
// 	}
// 	otto.Publish("ss/m/station/"+s.ID, string(json))
// }

// func (s *Station) Advertise(d time.Duration) {
// 	if s.ticker == nil {
// 		s.ticker = time.NewTicker(d)
// 		s.Announce()
// 	}

// 	s.quit = make(chan bool)
// 	go func() {
// 		for {
// 			select {
// 			case <-s.ticker.C:
// 				s.Announce()

// 			case <-s.quit:
// 				s.ticker.Stop()
// 				return
// 			}
// 		}
// 	}()
// }

// func (s Station) MarshalJSON() (j []byte, err error) {

// 	type Stat struct {
// 		ID        string             `json:"id"`
// 		LastHeard time.Time          `json:"last-heard"`
// 		Sensors   map[string]float64 `json:"sensors"`
// 	}

// 	stat := Stat{
// 		ID:        s.ID,
// 		LastHeard: s.LastHeard,
// 		Sensors:   make(map[string]float64),
// 	}

// 	for id, sens := range s.Sensors {
// 		data := sens.LastValue.(MsgData)

// 		v, err := strconv.ParseFloat(data.Value.(string), 64)
// 		if err != nil {
// 			log.Printf("ERROR StationJSON ParseFloat: %s %+v", data.Value.(string))
// 			v = -99.99
// 		}
// 		stat.Sensors[id] = v
// 	}

// 	j, err = json.Marshal(&stat)
// 	return j, err
// }
