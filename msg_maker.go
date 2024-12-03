package otto

import (
	"encoding/json"
	"math/rand/v2"
)

type MsgMaker interface {
	NewMsg() *Msg
}

type WeatherData struct {
	Tempc    float32 `json:"tempc"`
	Humidity float32 `json:"humidity"`
}

func (w *WeatherData) NewMsg() *Msg {
	w.Tempc = rand.Float32()
	w.Humidity = rand.Float32()

	msg := NewMsg()
	msg.Path = []string{"ss", "d", "station", "weather"}
	msg.Source = "MM"

	j, err := json.Marshal(w)
	if err != nil {
		l.Println("Error marshalling JSON data", err)
		return nil
	}

	msg.Message = j
	return msg
}
