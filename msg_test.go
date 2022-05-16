package main

import (
	"testing"	
)

func TestMsg(t *testing.T) {
	sm := NewStationManager()
	msg := Msg{
		Station: "localhost",
		Sensor:  "tester",
		Value: 72.5,
	}

	if msg.Station != "localhost" {
		t.Error("This is weird expected sanity")
	}

	sm.Recv(msg)

	st := sm.Get("localhost")
	if st == nil {

		t.Error("StationManager get expect (localhost) got (nil)")
		return
	}

	sens := st.GetSensor("tester")
	if sens == nil {
		t.Error("StationManager get sensor expect (tester) got (nil)")
		return
	}

	if len(sens.Values) != 1 {
		t.Errorf("sens.Count expected (1) got (%d)", len(sens.Values))
	}
}
