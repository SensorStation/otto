package data

import (
	"encoding/json"
	"fmt"

	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/message"
)

// DataManager is a map of Timeseries data that belongs to
// a specific station. The timeseries for each station are
// differentiated by the timeseries labels.
type DataManager struct {
	DataMap map[string]map[string]*Timeseries `json:"datamap"`
}

var (
	dataManager *DataManager
	l           *logger.Logger
)

func GetDataManager() *DataManager {
	if dataManager == nil {
		dataManager = NewDataManager()
	}
	if l == nil {
		l = logger.GetLogger()
	}
	return dataManager
}

// NewDataManager creates a new DataManager typically called
// by NewStation()
func NewDataManager() (dm *DataManager) {
	dm = &DataManager{
		DataMap: make(map[string]map[string]*Timeseries),
	}
	return dm
}

// Callback is the callback used by the DataManager to receive
// MQTT messages. TODO: move this call back to the stations because
// the stations will have a better understanding of the data they
// are subscribing to.
func (dm *DataManager) Callback(msg *message.Msg) {

	var m map[string]interface{}
	err := json.Unmarshal(msg.Data, &m)
	if err != nil {
		l.Error("Failed to unmarshal data: %s", "error", err)
		return
	}

	fmt.Printf("%+v\n", m)
}
