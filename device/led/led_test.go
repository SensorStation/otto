package led

import (
	"testing"

	"github.com/sensorstation/otto/device"
)

func TestLED(t *testing.T) {
	device.Mock(true)

	led := New("led", 5)
	if led.Name() != "led" {
		t.Errorf("led name got (%s) want (%s)", led.Name(), "led")
	}

}
