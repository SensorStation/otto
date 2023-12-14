package iote

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Msg holds a value and some type of meta data to be pass around in
// the system.
type Msg struct {
	Id       string      `json:"id"`
	Category string      `json:"category"`
	Station  string      `json:"station"` // mac addr
	Device   string      `json:"device"`
	Time     time.Time   `json:"time"`
	Value    interface{} `json:"value"`
}

func MsgFromMQTT(topic string, payload []byte) *Msg {

	// extract the station from the topic
	paths := strings.Split(topic, "/")

	if len(paths) < 3 {
		log.Println("[W] Unknown path: ", topic)
		return nil
	}

	// ss/msg/<source>/<sensor> <value>
	msg := &Msg{
		Station:  paths[1],
		Category: paths[2],
		Device:   paths[3],
		Time:     time.Now(),
		Value:    payload,
	}

	return msg
}

func (m Msg) String() string {
	var str string
	str = fmt.Sprintf("Time: %s, Category: %s, Station: %s, Device: %s = %q",
		m.Time.Format(time.RFC3339), m.Category, m.Station, m.Device, m.Value)
	return str
}
