package vh400

import (
	"testing"

	"github.com/sensorstation/otto/device"
)

func init() {
	device.Mock(true)
}

func TestVH400(t *testing.T) {
	v := New("vh400", 1)
	if v.Name() != "vh400" {
		t.Errorf("vh400 name expected (%s) got (%s)", "vh400", v.Name())
	}

}
