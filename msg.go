package otto

import (
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

	time.Time `json:"time"`
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

	// extract the station from the topic
	m.Path = strings.Split(topic, "/")
	m.Message = payload
	m.Time = time.Now()
	return &m, nil
}
