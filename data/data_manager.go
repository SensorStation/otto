package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/message"
	"github.com/sensorstation/otto/messanger"
)

// DataManager is a map of Timeseries data that belongs to
// a specific station. The timeseries for each station are
// differentiated by the timeseries labels.
type DataManager struct {
	dataMap map[string]map[string]*Timeseries `json:"datamap"`
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
		dataMap: make(map[string]map[string]*Timeseries),
	}
	mqtt := messanger.GetMQTT()
	if mqtt != nil {
		mqtt.Subscribe("ss/d/#", dm.Callback)
	}
	return dm
}

// Add will add data according to station and label
func (dm *DataManager) Add(station, label string, data any) {
	stmap, ex := dm.dataMap[station]
	if !ex {
		dm.dataMap[station] = make(map[string]*Timeseries)
		stmap = dm.dataMap[station]
	}
	ts, ex := stmap[label]
	if !ex {
		stmap[label] = NewTimeseries(station, label)
	}
	ts = stmap[label]
	ts.Add(data)
}

func (dm *DataManager) Dump(w io.Writer) {
	for _, st := range dm.dataMap {
		for _, ts := range st {
			fmt.Fprint(w, ts.String())
		}
	}
}

// Callback is the callback used by the DataManager to receive
// MQTT messages. TODO: move this call back to the stations because
// the stations will have a better understanding of the data they
// are subscribing to.
func (dm *DataManager) Callback(msg *message.Msg) {
	if msg.IsJSON() {
		m, err := msg.Map()
		if err != nil {
			return
		}

		for k, v := range m {
			dm.Add(msg.Station(), k, v)
		}
	}

	dm.Dump(os.Stdout)
	return
}

// ServeHTTP provides a REST interface to the config structure
func (dm DataManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		err := json.NewEncoder(w).Encode(dm)
		if err != nil {
			l.Error("Failed to encode data:", "error", err)
			http.Error(w, "failure", 401)
		}

	case "POST", "PUT":
		http.Error(w, "Not Yet Supported", 401)
	}
}
