package iote

import (
	"fmt"
	"strings"
	"time"
)

// Msg holds a value and some type of meta data to be pass around in
// the system.
type Msg struct {
	ID   int64       `json:"id"`
	Type string      `json:"type"`
	Data interface{} `json:"value"`

	time.Time `json:"time"`
}

type MsgData struct {
	Station string      `json:"station"` // mac addr
	Device  string      `json:"device"`
	Value   interface{} `json:"value"`
}

var (
	msgid int64 = 0
)

func getMsgID() int64 {
	msgid++
	return msgid
}

func MsgFromMQTT(topic string, payload []byte) (*Msg, error) {

	// extract the station from the topic
	paths := strings.Split(topic, "/")

	if len(paths) < 3 {
		err := fmt.Errorf("[E] Unknown path %s", topic)
		return nil, err
	}

	msg := &Msg{
		ID:   getMsgID(),
		Type: paths[1],
		Time: time.Now(),
	}
	data := MsgData{
		Value: string(payload),
	}

	switch msg.Type {
	case "m":
		data.Device = paths[2]
		data.Station = paths[3]

	case "d":
		data.Station = paths[2]
		data.Device = paths[3]
	}
	msg.Data = data
	return msg, nil
}

func (m Msg) String() string {
	var str string
	str = fmt.Sprintf("ID: %d, Time: %s, Type: %s ",
		m.ID, m.Time.Format(time.RFC3339), m.Type)

	switch m.Data.(type) {
	case MsgData:
		str += m.Data.(MsgData).String()

	case Station:
		str += m.Data.(Station).String()
	}
	return str
}

func (m MsgData) String() string {
	str := fmt.Sprintf("Station: %s, Device: %s = %q",
		m.Station, m.Device, m.Value)
	return str
}
