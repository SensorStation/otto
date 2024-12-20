package blasters

import (
	"encoding/json"
	"math/rand/v2"

	"github.com/sensorstation/otto"
)

// MsgMaker creates messages to be used by the mqtt blaster
// for smoke testing the messaging subsystem
type MsgMaker interface {
	NewMsg() *otto.Msg
}

// WeatherData is the content contained in the message used by the blaster
type WeatherData struct {
	Tempc    float32 `json:"tempc"`
	Humidity float32 `json:"humidity"`
}

// NewMsg will create a new message for testing
func (w *WeatherData) NewMsg() *otto.Msg {
	w.Tempc = rand.Float32()
	w.Humidity = rand.Float32()

	j, err := json.Marshal(w)
	if err != nil {
		otto.GetLogger().Error("Error marshalling JSON data", "error", err)
		return nil
	}

	path := "ss/d/station/weather"
	msg := otto.NewMsg(path, j, "weather-data")
	return msg
}
