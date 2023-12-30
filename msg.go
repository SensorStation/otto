package iote

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

func (m Msg) String() string {
	var str string
	str = fmt.Sprintf("ID: %d, Time: %s, Type: %s, ",
		m.ID, m.Time.Format(time.RFC3339), m.Type)
	str += m.Data.String()
	return str
}

func (m MsgStation) String() string {
	str := fmt.Sprintf("Station: %s, tempf: %f, humidity: %f, ",
		m.ID, m.Sensors["tempf"], m.Sensors["humidity"])

	for k, v := range m.Relays {
		vs := "off"
		if v {
			vs = "on"
		}
		str += k + ": " + vs
	}
	return str
}
