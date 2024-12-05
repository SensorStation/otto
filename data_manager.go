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

func (dm *DataManager) Callback(msg *Msg) {
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
