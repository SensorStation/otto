package otto

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func getMsg() (*Msg, time.Time) {
	now := time.Now()
	path := "ss/d/%s/test" 

	b := fmt.Sprintf("%d", 4)
	m := NewMsg(path, []byte(b))
	m.Source = "be:ef:ca:fe:01"
	m.Time = now

	return m, now
}

func TestStationMsg(t *testing.T) {
	topic := "ss/d/be:ef:ca:fe:01/station"
	omsg, _ := getMsg()
	omsg.Source = "be:ef:ca:fe:01"

	j, err := json.Marshal(omsg)
	if err != nil {
		t.Errorf("json marshal failed %+v", err)
		return
	}

	msg := NewMsg(topic, j)

	if msg == nil {
		t.Error("msg topic expected but is nil")
	}

}

func TestJSON(t *testing.T) {
	m, _ := getMsg()

	_, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Failed to marshal message %+v", m)
	}
}

func TestDataString(t *testing.T) {
	// m, now := getMsg()

	// formatted := fmt.Sprintf("ID: %d, Time: %s, Type: %s, Station: %s, tempf: %f, humidity: %f, ",
	// 	m.ID, now.Format(time.RFC3339), m.Type, m.Data.ID, m.Data.Sensors["tempf"], m.Data.Sensors["humidity"],
	// )

	// str := m.String()
	// if str != formatted {
	// 	t.Errorf("Data Formatted expected (%s) got (%s)", formatted, str)
	// }
}
