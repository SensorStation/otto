package blasters

import (
	"encoding/json"
	"math/rand/v2"
	"strings"

	"github.com/sensorstation/otto"
)

type MsgMaker interface {
	NewMsg() *otto.Msg
}

type WeatherData struct {
	Tempc    float32 `json:"tempc"`
	Humidity float32 `json:"humidity"`
}

func (w *WeatherData) NewMsg() *otto.Msg {
	w.Tempc = rand.Float32()
	w.Humidity = rand.Float32()

	path := strings.Join([]string{"ss", "d", "station", "weather"}, "/")
	data := []byte("MM")
	msg := otto.NewMsg(path, data)

	j, err := json.Marshal(w)
	if err != nil {
		otto.GetLogger().Error("Error marshalling JSON data", "error", err)
		return nil
	}

	msg.Message = j
	return msg
}
