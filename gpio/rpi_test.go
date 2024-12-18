package gpio

import (
	"encoding/json"
	"testing"
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
	var gpio GPIO
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
