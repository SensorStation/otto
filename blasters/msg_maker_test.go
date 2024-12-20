package blasters

import (
	"encoding/json"
	"testing"
)

func TestMsgMaker(t *testing.T) {
	wd := WeatherData{}
	msg := wd.NewMsg()
	exp := []string{"ss", "d", "station", "weather"}
	for i := 0; i < len(exp); i++ {
		if msg.Path[i] != exp[i] {
			t.Errorf("expected path index[%d] to be (%s) got (%s)", i, msg.Path[i], exp[i])
		}
	}

	if msg.Source != "weather-data" {
		t.Errorf("expected path to be (weather-data) got (%s)", msg.Source)
	}

	var j WeatherData
	err := json.Unmarshal(msg.Message, &j)
	if err != nil {
		t.Errorf("failed to unmarshal message %s", err)
	}

	if j.Tempc == 0 {
		t.Errorf("tempc is not random got 0")
	}

	if j.Humidity == 0 {
		t.Errorf("humidity is not random got 0")
	}

}
