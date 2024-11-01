package otto

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

	t := now.Format(time.RFC3339)
	m := Msg{
		ID:   1,
		Type: "d",
		Time: t,
		Data: s,
	}
	return m, now
}

func TestStationMsg(t *testing.T) {
	topic := "ss/d/be:ef:ca:fe:01/station"
	omsg, _ := getMsg()

	j, err := json.Marshal(omsg)
	if err != nil {
		t.Errorf("json marshal failed %+v", err)
		return
	}

	// t.Logf("j: %s\n", string(j))
	msg, err := MsgFromMQTT(topic, j)
	if err != nil {
		t.Errorf("Extracting message from MQTT %+v", err)
		return
	}

	t.Logf("MSG: %+v", msg)
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
