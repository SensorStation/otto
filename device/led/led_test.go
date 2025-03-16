package led

import (
	"testing"

	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/messanger"
)

func TestLED(t *testing.T) {
	device.Mock(true)

	led := New("led", 5)
	if led.Name() != "led" {
		t.Errorf("led name got (%s) want (%s)", led.Name(), "led")
	}

	msg := messanger.New(led.Topic, []byte("on"), "test")
	led.Callback(msg)

	v, err := led.Value()
	if err != nil {
		t.Fatalf("led.Value() got error %v", err)
	}
	if v != 1 {
		t.Errorf("led expected (1) got (%d)", v)
	}

	msg = messanger.New(led.Topic, []byte("off"), "test")
	led.Callback(msg)

	v, err = led.Value()
	if err != nil {
		t.Fatalf("led.Value() got error %v", err)
	}
	if v != 0 {
		t.Errorf("led expected (0) got (%d)", v)
	}

}
