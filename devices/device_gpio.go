package devices

import (
	"time"

	"github.com/sensorstation/otto"
	"github.com/warthog618/go-gpiocdev"
)

type GPIODevice struct {
	*Dev
	*Pin
	EvtQ chan gpiocdev.LineEvent
}

func (d *GPIODevice) On() {
	d.Set(1)
}

func (d *GPIODevice) Off() {
	d.Set(0)
}

func (d *GPIODevice) Set(v int) {
	d.Pin.Set(v)
	val := "off"
	if v > 0 {
		val = "on"
	}

	m := otto.GetMQTT()
	for _, p := range d.pubs {
		m.Publish(p, val)
	}
}

func GPIOOut(name string, pin int) *GPIODevice {
	g := GetGPIO()
	d := &GPIODevice{
		Dev: &Dev{
			name: name,
		},
		Pin: g.Pin(name, pin, gpiocdev.AsOutput(0)),
	}

	topic := "ss/c/" + stationName + "/" + name
	d.subs = append(d.subs, topic)

	m := otto.GetMQTT()
	m.Subscribe(topic, d)

	d.pubs = append(d.pubs, "ss/d/"+stationName+"/"+name)
	return d
}

func GPIOIn(name string, pin int) *GPIODevice {
	g := GetGPIO()
	d := &GPIODevice{
		Dev: &Dev{
			name: name,
		},
		EvtQ: make(chan gpiocdev.LineEvent),
	}
	d.Pin = g.Pin(name, pin,
		gpiocdev.WithPullUp,
		gpiocdev.WithFallingEdge,
		gpiocdev.WithDebounce(10*time.Millisecond),
		gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
			d.EvtQ <- evt
		}))

	d.pubs = append(d.pubs, "ss/c/"+stationName+"/"+name)
	return d
}
