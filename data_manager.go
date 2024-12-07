package otto

import (
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

func (dm *DataManager) SubCallback(topic string, message []byte) {

	// convert the topic and data into a *Msg
	msg := NewMsg(topic, message)
	if len(msg.Path) < 3 {
		l.Error("DataManager: Malformed MQTT ", "path", msg.Path)
		return
	}

	// Change this to a map[string]string or map[string]interface{}
	st := stations.Update(msg)

	var m map[string]interface{}
	err := json.Unmarshal(msg.Message, &m)
	if err != nil {
		l.Error("Failed to unmarshal message ", "error", err)
		return
	}
	for k, v := range m {
		st.Insert(k, v)
	}
}
