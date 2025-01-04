package message

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sensorstation/otto/utils"
)

type Message interface {
}

// Msg holds a value and some type of meta data to be pass around in
// the system.
type Msg struct {
	ID     int64    `json:"id"`
	Topic  string   `json:"topic"`
	Path   []string `json:"path"`
	Args   []string `json:"args"`
	Data   []byte   `json:"msg"`
	Source string   `json:"source"`

	Timestamp time.Duration `json:"timestamp"`
}

var (
	msgid    int64 = 0
	msgSaver *MsgSaver
)

func getMsgID() int64 {
	msgid++
	return msgid
}

func New(topic string, data []byte, source string) *Msg {
	msg := &Msg{
		ID:        getMsgID(),
		Topic:     topic,
		Path:      strings.Split(topic, "/"),
		Data:      data,
		Timestamp: utils.Timestamp(),
		Source:    source,
	}

	if msgSaver != nil && msgSaver.Saving {
		msgSaver.Messages = append(msgSaver.Messages, msg)
	}
	return msg
}

func (msg *Msg) Station() string {
	if len(msg.Path) < 3 {
		return ""
	}
	return msg.Path[3]
}

func (msg *Msg) Last() string {
	l := len(msg.Path)
	return msg.Path[l-1]
}

func (msg *Msg) Byte() []byte {
	return msg.Data
}

func (msg *Msg) String() string {
	return string(msg.Data)
}

func (msg *Msg) Float64() float64 {
	var f float64
	fmt.Sscanf(msg.String(), "%f", &f)
	return f
}

func (msg *Msg) IsJSON() bool {
	return json.Valid(msg.Data)
}

func (msg *Msg) Map() (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(msg.Data, &m)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal data: %s", err)
	}
	return m, nil
}

func (msg *Msg) Dump() string {
	str := fmt.Sprintf("  ID: %d\n", msg.ID)
	str += fmt.Sprintf("Path: %q\n", msg.Path)
	str += fmt.Sprintf("Args: %q\n", msg.Args)
	str += fmt.Sprintf(" Src: %s\n", msg.Source)
	str += fmt.Sprintf("Time: %s\n", msg.Timestamp)
	str += fmt.Sprintf("Data: %s\n", string(msg.Data))

	return str
}

type MsgSaver struct {
	Messages []*Msg `json:"saved-messages"`
	Saving   bool   `json:"saving"`
}

func GetMsgSaver() *MsgSaver {
	if msgSaver == nil {
		msgSaver = &MsgSaver{}
	}
	return msgSaver
}

func (ms *MsgSaver) StartSaving() {
	ms.Saving = true
}

func (ms *MsgSaver) StopSaving() {
	ms.Saving = false
}

func (ms *MsgSaver) Dump() {
	for _, msg := range ms.Messages {
		println(msg.Dump())
		println("----------------------------------------------")
	}
}

// ServeHTTP will respond to the writer with 'Pong'
func (ms *MsgSaver) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms)
}
