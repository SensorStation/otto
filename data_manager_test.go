package otto

import (
	"fmt"
	"testing"

	"github.com/sensorstation/otto/message"
)

func TestNewDataManager(t *testing.T) {
	dm := NewDataManager()
	if len(dm.DataMap) != 0 {
		t.Errorf("Datamanger map not empty expected(0) got (%d)", len(dm.DataMap))
	}
}

func TestDataManagerCallback(t *testing.T) {
	data := []byte(fmt.Sprintf(`{ "int": 10 }`))
	path := "ss/d/station1/test"

	dm := NewDataManager()
	msg := message.NewMsg(path, data, "data-manager-test")
	dm.Callback(msg)

	sm := GetStationManager()
	st := sm.Get("station1")
	if len(st.DataManager.DataMap) != 1 {
		t.Errorf("failed to get count == 1 data from station1")
	}
}
