package iote

import (
	"testing"
)

func TestStation(t *testing.T) {
	localip := "127.0.0.1"
	st := NewStation(localip)

	if st.ID != localip {
		t.Errorf("IP expecting (%s) got (%s)", localip, st.ID)
	}
}

func TestStationManager(t *testing.T) {
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

	if sm.Count() != len(ids) {
		t.Errorf("Station Manager count got (%d) expected (%d)",
			len(ids), sm.Count())
	}

	for _, id := range ids {
		st := sm.Get(id)
		if st == nil {
			t.Errorf("Get station expected (%s) got nothing", id)
		}
	}

}
