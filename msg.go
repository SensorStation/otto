package otto

import (
	"fmt"
	"strings"
	"time"
)

// Msg holds a value and some type of meta data to be pass around in
// the system.
type Msg struct {
	ID      int64    `json:"id"`
	Path    []string `json:"path"`
	Args    []string `json:"args"`
	Message []byte   `json:"msg"`
	Source  string   `json:"source"`
	Time    string   `json:"time"`
}

var (
	msgid int64 = 0
)

func getMsgID() int64 {
	msgid++
	return msgid
}

// MsgFromMQTT will parse the topic and pass the payload
// to the correct station for the given value.
func MsgFromMQTT(topic string, payload []byte) (*Msg, error) {

	var m Msg = Msg{}
	m.ID = getMsgID()

	fmt.Printf("TOPIC: %s\n", topic)

	// extract the station from the topic
	m.Path = strings.Split(topic, "/")
	if len(m.Path) < 3 {
		err := fmt.Errorf("[E] Unknown path %s", topic)
		return nil, err
	}

	m.Message = payload
	return &m, nil
}

// String will stringify the payload and topic from MQTT
func (m *Msg) OString() string {
	now := time.Now()

	formatted := fmt.Sprintf("ID: %d, Time: %s, Station: %s",
		m.ID, now.Format(time.RFC3339), m.ID)
	return formatted
}
