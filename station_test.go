package main

import (
	"testing"
)

func TestStation(t *testing.T) {
	localip := "127.0.0.1"
	st := NewStation(localip)

	if st.ID != localip {
		t.Errorf("IP expecting (%s) got (%s)", localip, st.ID)
	}

	if st.Length() != 0 {
		t.Errorf("Length expecting (0) got (%d)", st.Length())
	}

	st.AddData("test1", 22.3)
	if st.Length() != 1 {
		t.Errorf("Length expecting (1) got (%d)", st.Length())
	}

	for i := 0; i < 100; i++ {
		st.AddData("test2",  float64(i))
	}

	if st.Length() != 2 {
		t.Errorf("Length expecting (2) got (%d)", st.Length())
	}

	sens := st.GetSensor("test2")
	if sens == nil {
		t.Error("Expected sensor (test2) got (nil)")
	}

	cnt := len(sens.Values)
	if cnt != 100 {
		t.Errorf("DataCount test2 expected (100) got (%d)", cnt)
	}

	for i := 1; i < 100; i++ {
		ts := sens.Values[i]
		if ts.Val != float64(i) {
			t.Errorf("GetSensor test2 expected (%3.2f) got (%3.2f)", float64(i), ts.Val)
		}
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
