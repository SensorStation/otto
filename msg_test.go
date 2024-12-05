package otto

import (
	"encoding/json"
	"testing"
	"time"
)

func getMsg() (*Msg, time.Time) {
	now := time.Now()
	m := NewMsg()
	m.Source = "be:ef:ca:fe:01"
	m.Time = now
	// s.Timeseries["tempf"].Data = append(s.Timeseries["tempf"].Data, 97.8)
	// s.Timesereis["humidity"].Data = append(s.Timeseries["humidity"].Data, 99.3)

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

	dm := GetDataManager()
	msg := dm.GetMsg(topic, j)

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
