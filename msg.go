package iote

import (
	"fmt"
	"time"
)

// Msg holds a value and some type of meta data to be pass around in
// the system.
type Msg struct {
	Id       string    `json:"id"`
	Category string    `json:"category"`
	Station  string    `json:"station"` // mac addr
	Device   string    `json:"device"`
	Time     time.Time `json:"time"`
}

func (m Msg) String() string {
	var str string
	str = fmt.Sprintf("Time: %s, Category: %s, Station: %s, Device: %s",
		m.Time.Format(time.RFC3339),
		m.Category, m.Station, m.Device)
	return str
}

type MsgFloat64 struct {
	Msg
	Value float64 `json:"value"`
}

func (m MsgFloat64) String() string {
	str := m.Msg.String()
	str += fmt.Sprintf(" = %f", m.Value)
	return str
}
