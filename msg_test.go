package iote

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func getMsg() (Msg, time.Time) {
	now := time.Now()
	s := MsgStation{
		ID:      "be:ef:ca:fe:01",
		Sensors: make(map[string]float64),
		Relays:  make(map[string]bool),
	}
	s.Sensors["tempf"] = 97.8
	s.Sensors["humidity"] = 99.3

	m := Msg{
		ID:   1,
		Type: "d",
		Time: now,
		Data: s,
	}
	return m, now
}

func TestJSON(t *testing.T) {
	m, _ := getMsg()

	_, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Failed to marshal message %+v", m)
	}
}

func TestDataString(t *testing.T) {
	m, now := getMsg()

	formatted := fmt.Sprintf("ID: %d, Time: %s, Type: %s, Station: %s, tempf: %f, humidity: %f, ",
		m.ID, now.Format(time.RFC3339), m.Type, m.Data.ID, m.Data.Sensors["tempf"], m.Data.Sensors["humidity"],
	)

	str := m.String()
	if str != formatted {
		t.Errorf("Data Formatted expected (%s) got (%s)", formatted, str)
	}
}
