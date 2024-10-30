package gpio

import (
	"testing"
)

func TestPinMap(t *testing.T) {
	pin := pm.Find(4)
	if pin == nil {
		t.Errorf("Expected a pin but got nothing")
	}
}

func TestPin(t *testing.T) {

}
