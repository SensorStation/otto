package iote

import "time"

type Control struct {
	ID      string
	Command string
	Value   interface{}

	time.Time
}
