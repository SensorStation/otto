package station

import (
	"encoding/json"
	"testing"
	"time"
)

func StationCreation(count int) []string {
	ids := []string{
		"127.0.0.1",
		"127.0.0.2",
		"127.0.0.3",
		"127.0.0.4",
		"127.0.0.5",
	}

	sm := NewStationManager()
	for _, id := range ids {
		sm.Add(id)
	}
	return ids
}

func TestStation(t *testing.T) {
	localip := "127.0.0.1"
	st := NewStation(localip)

	if st.ID != localip {
		t.Errorf("IP expecting (%s) got (%s)", localip, st.ID)
	}
}

func TestStationManager(t *testing.T) {

	count := 5

	sm := NewStationManager()
	sids := StationCreation(count)
	for _, id := range sids {
		sm.Add(id)
	}

	if sm.Count() != len(sids) {
		t.Errorf("Station Manager count got (%d) expected (%d)",
			len(sids), sm.Count())
	}

	for _, id := range sids {
		st := sm.Get(id)
		if st.ID != id {
			t.Errorf("Get station expected (%s) got nothing", id)
		}
	}

}

func TestStationJSON(t *testing.T) {
	sens := make(map[string]float64)
	sens["tempf"] = 89.43
	sens["humidity"] = 99.00

	relays := make(map[string]bool)
	relays["relay0"] = false
	relays["relay1"] = true

	st := &Station{
		ID:        "aa:bb:cc:dd:ee:11",
		LastHeard: time.Now(),
	}

	j, err := json.Marshal(st)
	if err != nil {
		t.Errorf("Marshal Station failed: %+v", err)
		return
	}

	var station Station
	err = json.Unmarshal(j, &station)
	if err != nil {
		t.Errorf("Unmarshal Station failed: %+v", err)
		return
	}
}
