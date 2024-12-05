package otto

import (
	"strings"
	"time"

	"encoding/json"
)

type DataManager struct {
	Timeseries map[string][]*Data
}

func NewDataManager() (sm *DataManager) {
	sm = &DataManager{
		Timeseries: make(map[string][]*Data),
	}
	return sm
}

func (dm *DataManager) GetMsg(topic string, data []byte) *Msg {

	m := NewMsg()

	// extract the station from the topic
	m.Path = strings.Split(topic, "/")
	m.Message = data
	m.Time = time.Now()
	return m

}

func (dm *DataManager) SubCallback(topic string, message []byte) {

	// convert the topic and data into a *Msg
	msg := dm.GetMsg(topic, message)
	if len(msg.Path) < 3 {
		l.Printf("DataManager: Malformed MQTT path: %q\n", msg.Path)
		return
	}

	// Change this to a map[string]string or map[string]interface{}
	st := stations.Update(msg)

	var m map[string]interface{}
	err := json.Unmarshal(msg.Message, &m)
	if err != nil {
		l.Println("Failed to unmarshal message ", err)
		return
	}
	for k, v := range m {
		st.Insert(k, v)
	}
}
