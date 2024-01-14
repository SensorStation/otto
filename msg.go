package otto

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"
)

// Msg holds a value and some type of meta data to be pass around in
// the system.
type Msg struct {
	ID   int64      `json:"id"`
	Type string     `json:"type"`
	Data MsgStation `json:"station"`

	time.Time `json:"time"`
}

type MsgStation struct {
	ID      string             `json:"id"`
	Sensors map[string]float64 `json:"sensors"`
	Relays  map[string]bool    `json:"relays"`
}

var (
	msgid int64 = 0
)

func getMsgID() int64 {
	msgid++
	return msgid
}

func MsgFromMQTT(topic string, payload []byte) (m *Msg, err error) {

	// extract the station from the topic
	paths := strings.Split(topic, "/")
	if len(paths) < 3 {
		err := fmt.Errorf("[E] Unknown path %s", topic)
		return nil, err
	}

	m = &Msg{
		ID:   getMsgID(),
		Type: paths[1],
		Time: time.Now(),
	}

	var data MsgStation
	err = json.Unmarshal(payload, &data)
	if err != nil {
		return nil, err
	}
	m.Data = data
	return m, nil
}

func (m *Msg) String() string {
	now := time.Now()

	formatted := fmt.Sprintf("ID: %d, Time: %s, Type: %s, Station: %s, tempf: %f, humidity: %f, ",
		m.ID, now.Format(time.RFC3339), m.Type, m.Data.ID, m.Data.Sensors["tempf"], m.Data.Sensors["humidity"])
	return formatted
}
