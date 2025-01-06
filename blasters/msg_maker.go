package blasters

import (
	"encoding/json"
	"log/slog"
	"math/rand/v2"

	"github.com/sensorstation/otto/message"
)

// MsgMaker creates messages to be used by the mqtt blaster
// for smoke testing the messaging subsystem
type MsgMaker interface {
	NewMsg() *message.Msg
}

// WeatherData is the content contained in the message used by the blaster
type WeatherData struct {
	Tempc    float32 `json:"tempc"`
	Humidity float32 `json:"humidity"`
}

// NewMsg will create a new message for testing
func (w *WeatherData) NewMsg(topic string) *message.Msg {
	w.Tempc = rand.Float32()
	w.Humidity = rand.Float32()

	j, err := json.Marshal(w)
	if err != nil {
		slog.Error("Error marshalling JSON data", "error", err)
		return nil
	}

	msg := message.New(topic, j, "weather-data")
	return msg
}
