package iote

import (
	"fmt"
	"strings"
	"time"
)

// Msg holds a value and some type of meta data to be pass around in
// the system.
type Msg struct {
	Id       string      `json:"id"`
	Topic    string      `json:"topic"`
	Category string      `json:"category"`
	Station  string      `json:"station"` // mac addr
	Device   string      `json:"device"`
	Time     time.Time   `json:"time"`
	Value    interface{} `json:"value"`
}

func MsgFromMQTT(topic string, payload []byte) (*Msg, error) {

	// extract the station from the topic
	paths := strings.Split(topic, "/")

	if len(paths) < 3 {
		err := fmt.Errorf("[E] Unknown path %s", topic)
		return nil, err
	}

	// ss/C/<source>/<sensor> <value>
	msg := &Msg{
		Topic:    topic,
		Category: paths[1],
		Station:  paths[2],
		Device:   paths[3],
		Time:     time.Now(),
	}
	msg.Value = payload
	Stations.Update(msg)

	return msg, nil
}

func (m Msg) String() string {
	var str string
	str = fmt.Sprintf("Time: %s, Category: %s, Station: %s, Device: %s = %q",
		m.Time.Format(time.RFC3339), m.Category, m.Station, m.Device, m.Value)
	return str
}
