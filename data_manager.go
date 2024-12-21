package otto

import (
	"encoding/json"
	"fmt"

	"github.com/sensorstation/otto/message"
)

// DataManager is a map of Timeseries data that belongs to
// a specific station. The timeseries for each station are
// differentiated by the timeseries labels.
type DataManager struct {
	DataMap map[string]*Timeseries `json:"datamap"`
}

// NewDataManager creates a new DataManager typically called
// by NewStation()
func NewDataManager() (dm *DataManager) {
	dm = &DataManager{
		DataMap: make(map[string]*Timeseries),
	}
	return dm
}

// SubCallback is the callback used by the DataManager to receive
// MQTT messages. TODO: move this call back to the stations because
// the stations will have a better understanding of the data they
// are subscribing to.
func (dm *DataManager) SubCallback(msg *message.Msg) {
	// Change this to a map[string]string or map[string]interface{}
	stations := GetStationManager()
	st := stations.Update(msg)

	var m map[string]interface{}
	err := json.Unmarshal(msg.Data, &m)
	if err != nil {
		l.Error("Failed to unmarshal message ", "error", err)
		return
	}
	for k, v := range m {
		st.Insert(k, v)
	}
	fmt.Printf("ST: %+v\n", st.DataManager)
}
