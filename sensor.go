package otto

import (
	"fmt"
	"time"
)

type Sensor struct {
	ID        string `json:"id"`
	LastValue interface{}
	LastHeard time.Time `json:"-"`
}

func (s *Sensor) Update(value interface{}, t time.Time) {
	s.LastValue = value
	s.LastHeard = t
}
func (s *Sensor) Callback(msg *Msg) {
	fmt.Printf("sensor: %+v", msg)
}
