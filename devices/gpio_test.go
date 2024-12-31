package devices

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/sensorstation/otto/message"
	"github.com/warthog618/go-gpiocdev"
)

var gpioStr = `
{
    "chipname":"gpiochip4",
    "pins": {
        "6": {
            "name": "led",
            "offset": 6,
            "value": 0,
            "mode": 0
        }
    }
}
`

func TestFromJSON(t *testing.T) {
	gpio := GetGPIO()
	gpio.Mock = true
	if err := json.Unmarshal([]byte(gpioStr), &gpio); err != nil {
		t.Error(err)
	}

	if gpio.Chipname != "gpiochip4" {
		t.Errorf("expected chipname (gpiochip4) got (%s)", gpio.Chipname)
	}

	if err := gpio.Init(); err != nil {
		t.Error(err)
	}
}

func TestPin(t *testing.T) {
	g := GetGPIO()
	g.Mock = true

	p := &Pin{
		offset: 1,
		mock:   true,
	}
	p.Opts = append(p.Opts, gpiocdev.AsOutput(1))
	err := p.Init()
	if err != nil {
		t.Error("Failed to initialize the pin", err)
	}

	v, err := p.Get()
	if err != nil {
		t.Errorf("Failed to read the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 1 {
		t.Errorf("Expected pin[%d] to be (1) got (%d)", p.Offset(), v)
	}

	err = p.Off()
	v, err = p.Get()
	if err != nil {
		t.Errorf("Failed to read the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 0 {
		t.Errorf("Expected pin[%d] to be (1) got (%d)", p.Offset(), v)
	}

	err = p.Set(0)
	v, err = p.Get()
	if err != nil {
		t.Errorf("Failed to read the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 0 {
		t.Errorf("Expected pin[%d] to be (1) got (%d)", p.Offset(), v)
	}

	err = p.On()
	v, err = p.Get()
	if err != nil {
		t.Errorf("Failed to read the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 1 {
		t.Errorf("Expected pin[%d] to be (1) got (%d)", p.Offset(), v)
	}

	// Now reconfigure the pin to be an output pin
	err = p.Reconfigure(gpiocdev.AsInput)
	if err != nil {
		t.Error("Failed to reconfigure pin AsInput", err)
		return
	}

	v, err = p.Get()
	if err != nil {
		t.Errorf("Getting the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 1 {
		t.Errorf("Error retrieving the value of pin[%d] expected(1) got (%d)", p.Offset(), v)
	}

	// Fake a hwset
	// p.SetValue(1)
	l := p.Line.(*MockLine)
	l.MockHWInput(0)
	v, err = p.Get()
	if err != nil {
		t.Errorf("Getting the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 0 {
		t.Errorf("Error retrieving the value of pin[%d] expected(1) got (%d)", p.Offset(), v)
	}

	err = p.Toggle()
	if err != nil {
		t.Error("Errored on toggle", err)
	}

	v, err = p.Get()
	if err != nil {
		t.Errorf("Getting the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 1 {
		t.Errorf("Error retrieving the value of pin[%d] expected(1) got (%d)", p.Offset(), v)
	}

	msg := message.New("ss/c/station/pin", []byte("off"), "test-pin")
	p.Callback(msg)
	v, err = p.Get()
	if err != nil {
		t.Errorf("Getting the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 0 {
		t.Errorf("Error retrieving the value of pin[%d] expected(1) got (%d)", p.Offset(), v)
	}

	msg = message.New("ss/c/station/pin", []byte("on"), "test-pin")
	p.Callback(msg)
	v, err = p.Get()
	if err != nil {
		t.Errorf("Getting the value of pin[%d] - %s", p.Offset(), err)
	}

	if v != 1 {
		t.Errorf("Error retrieving the value of pin[%d] expected(1) got (%d)", p.Offset(), v)
	}

	str := p.String()
	str = strings.TrimSpace(str)
	if !strings.Contains(str, "1: 1") {
		t.Errorf("expected string (1: 1) got (%s)", str)
	}

	p.Close()

	// clear up GPIO for future tests
	g = GetGPIO()
	g.Shutdown()
}

func TestGPIO(t *testing.T) {
	g := GetGPIO()
	g.Mock = true
	g.Init()

	if len(g.Pins) > 0 {
		t.Errorf("Pins are defined expected (0) got (%d)", len(g.Pins))
	}

	npins := 30
	for i := 0; i < npins; i++ {
		if i%2 == 0 {
			g.Pin(strconv.Itoa(i), i, gpiocdev.AsInput)
		} else {
			g.Pin(strconv.Itoa(i), i, gpiocdev.AsOutput(1))
		}
	}

	if len(g.Pins) != 30 {
		t.Errorf("Expected pin count (%d) got (%d)", npins, len(g.Pins))
	}

	str := g.String()
	if str == "" {
		t.Error("expected gpio.String() output but got nothing")
	}

	g.Shutdown()
	if g.Pins != nil {
		t.Errorf("Expected gpio.Pins to be nil but got something")
	}
}