package data

import (
	"fmt"
	"testing"

	"github.com/sensorstation/otto/message"
)

func TestNewDataManager(t *testing.T) {
	dm := NewDataManager()
	if len(dm.dataMap) != 0 {
		t.Errorf("Datamanger map not empty expected(0) got (%d)", len(dm.dataMap))
	}
}

func TestDataManagerCallback(t *testing.T) {
	data := []byte(fmt.Sprintf(`{ "int": 10 }`))
	path := "ss/d/station1/test"

	dm := NewDataManager()
	msg := message.New(path, data, "data-manager-test")
	dm.Callback(msg)
}
