package otto

import (
	"encoding/json"
	"math/rand/v2"
	"strings"
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

	path := strings.Join([]string{"ss", "d", "station", "weather"}, "/")
	data := []byte("MM")
	msg := NewMsg(path, data)

	j, err := json.Marshal(w)
	if err != nil {
		l.Error("Error marshalling JSON data", "error", err)
		return nil
	}

	msg.Message = j
	return msg
}
