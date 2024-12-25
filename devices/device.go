package devices

import "time"

type Device interface {
	Name() string
	SetName(name string)
	Pubs() []string
	AddPub(topic string)
}

type Dev struct {
	name string
	pubs []string
	subs []string

	period time.Duration
}

func (d *Dev) Name() string {
	return d.name
}

func (d *Dev) SetName(name string) {
	d.name = name
}

func (d *Dev) Pubs() []string {
	return d.pubs
}

func (d *Dev) AddPub(topic string) {
	d.pubs = append(d.pubs, topic)
}

func (d *Dev) Period() time.Duration {
	return d.period
}
