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

	time.Time `json:"time"`
}

var (
	msgid int64 = 0
)

func getMsgID() int64 {
	msgid++
	return msgid
}

func NewMsg() *Msg {
	msg := &Msg{
		ID: getMsgID(),
	}
	return msg
}

// MsgFromMQTT will parse the topic and pass the payload
// to the correct station for the given value.
func MsgFromMQTT(topic string, payload []byte) (*Msg, error) {

	m := NewMsg()

	// extract the station from the topic
	m.Path = strings.Split(topic, "/")
	m.Message = payload
	m.Time = time.Now()
	return m, nil
}

func (msg *Msg) Byte() []byte {
	return msg.Message
}

func (msg *Msg) String() string {
	str := fmt.Sprintf("  ID: %d\n", msg.ID)
	str += fmt.Sprintf("Path: %q\n", msg.Path)
	str += fmt.Sprintf("Args: %q\n", msg.Args)
	str += fmt.Sprintf(" Msg: %s\n", string(msg.Message))
	str += fmt.Sprintf(" Src: %s\n", msg.Source)
	str += fmt.Sprintf("Time: %s\n", msg.Time)
	return str
}
