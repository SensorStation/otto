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

	Period time.Duration
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

func (d *Dev) Subs() []string {
	return d.pubs
}

func (d *Dev) AddPub(topic string) {
	d.pubs = append(d.pubs, topic)
}
