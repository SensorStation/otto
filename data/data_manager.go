package data

import (
	"fmt"

	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/message"
	"github.com/sensorstation/otto/messanger"
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
	mqtt := messanger.GetMQTT()
	if mqtt != nil {
		mqtt.Subscribe("ss/d/#", dm.Callback)
	}
	return dm
}

// Callback is the callback used by the DataManager to receive
// MQTT messages. TODO: move this call back to the stations because
// the stations will have a better understanding of the data they
// are subscribing to.
func (dm *DataManager) Callback(msg *message.Msg) error {
	fmt.Printf("M: %+v\n", msg)
	if msg.IsJSON() {
		m, err := msg.Map()
		if err != nil {
			return err
		}
		fmt.Printf("MAP: %+v\n", m)
	}

	return nil
}
