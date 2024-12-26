package devices

import (
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type Mode int

const (
	ModeNone Mode = iota
	ModeInput
	ModeOutput
	ModePWM
)

type Device struct {
	Name string
	Pubs []string
	Subs []string
	Mode

	Period time.Duration
	EvtQ   chan gpiocdev.LineEvent
}
