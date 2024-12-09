package otto

import (
	"fmt"
	"strings"
	"time"
)

type Message interface {
	GetMsg() *Msg
}

// Msg holds a value and some type of meta data to be pass around in
// the system.
type Msg struct {
	ID      int64    `json:"id"`
	Path    []string `json:"path"`
	Args    []string `json:"args"`
	Message []byte	 `json:"msg"`
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

func NewMsg(topic string, data []byte) *Msg {
	msg := &Msg{
		ID:      getMsgID(),
		Path:    strings.Split(topic, "/"),
		Message: data,
		Time:    time.Now(),
	}

	return msg
}

func (msg *Msg) Byte() []byte {
	return msg.Message
}

func (msg *Msg) String() string {
	return string(msg.Message)
}

func (msg *Msg) Dump() string {
	str := fmt.Sprintf("  ID: %d\n", msg.ID)
	str += fmt.Sprintf("Path: %q\n", msg.Path)
	str += fmt.Sprintf("Args: %q\n", msg.Args)
	str += fmt.Sprintf(" Msg: %s\n", string(msg.Message))
	str += fmt.Sprintf(" Src: %s\n", msg.Source)
	str += fmt.Sprintf("Time: %s\n", msg.Time)
	return str
}
