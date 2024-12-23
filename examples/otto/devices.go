package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/gpio"
	"github.com/sensorstation/otto/message"
	"github.com/warthog618/go-gpiocdev"
)

type Device struct {
	*gpio.Pin

	Subs []string
	Pubs []string
}

func (d *Device) On() {
	d.Set(1)
}

func (d *Device) Off() {
	d.Set(0)
}

func (d *Device) Set(v int) {
	d.Pin.Set(v)
	val := "off"
	if v > 0 {
		val = "on"
	}

	m := otto.GetMQTT()
	for _, p := range d.Pubs {
		m.Publish(p, val)
	}
}

type DeviceManager struct {
	devices map[string]*Device
}

var (
	stationName string = "station"
	devices     DeviceManager
	evtQ        chan gpiocdev.LineEvent
)

func init() {
	devices.devices = make(map[string]*Device)
	evtQ = make(chan gpiocdev.LineEvent)
}

func (dm *DeviceManager) Add(d *Device) {
	if dm.devices == nil {
		dm.devices = make(map[string]*Device)
	}
	dm.devices[d.Name] = d
}

func (dm *DeviceManager) Get(name string) *Device {
	d, ex := dm.devices[name]
	if !ex {
		l.Error("device does not exist", "device", name)
		return nil
	}
	return d
}

func (dm *DeviceManager) FindPin(offset int) *Device {
	for _, d := range dm.devices {
		if d.Offset() == offset {
			return d
		}
	}
	return nil
}

func GPIOOut(name string, pin int) *Device {
	g := gpio.GetGPIO()
	d := &Device{
		Pin: g.Pin(name, pin, gpiocdev.AsOutput(0)),
	}

	topic := "ss/c/" + stationName + "/" + name
	d.Subs = append(d.Subs, topic)

	m := otto.GetMQTT()
	m.Subscribe(topic, d)

	d.Pubs = append(d.Pubs, "ss/d/"+stationName+"/"+name)
	return d
}

func (dm *DeviceManager) SubCallback(msg *message.Msg) {

	if len(msg.Path) < 4 || msg.Path[3] != "button" {
		l.Error("callback from unwanted device", "device", msg.Path[4])
		return
	}
	led, ex := dm.devices["led"]
	if !ex {
		l.Error("failed to find led")
		return
	}

	relay, ex := dm.devices["relay"]
	if !ex {
		l.Error("failed to find relay")
		return
	}

	val := string(msg.Data)
	switch val {
	case "on":
	case "1":
		led.On()
		relay.On()

	case "off":
	case "0":
		led.Off()
		relay.Off()

	default:
		l.Error("Dont know what to do with", "value", val)
	}
}

func GPIOIn(name string, pin int) *Device {
	g := gpio.GetGPIO()
	d := &Device{
		Pin: g.Pin(name, pin,
			// gpiocdev.WithPullUp,
			gpiocdev.WithFallingEdge,
			gpiocdev.WithDebounce(10*time.Millisecond),
			gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
				evtQ <- evt
			})),
	}
	d.Pubs = append(d.Pubs, "ss/c/"+stationName+"/"+name)
	return d
}

func ButtonHandler() {
	for {
		select {
		case evt := <-evtQ:
			dev := devices.FindPin(evt.Offset)
			if dev == nil {
				// something went terribly wrong
				l.Error("Failed to find device for pin", "offset", evt.Offset)
				return
			}

			evtype := "falling"
			switch evt.Type {
			case gpiocdev.LineEventFallingEdge:
				evtype = "falling"

			case gpiocdev.LineEventRisingEdge:
				evtype = "raising"

			default:
				l.Warn("Unknown event type ", "type", evt.Type)
				continue
			}

			l.Info("GPIO edge", "device", dev.Name, "direction", evtype,
				"seqno", evt.Seqno, "lineseq", evt.LineSeqno)

			// Get the value of the relay then toggle it, that is the value we will send
			relay := devices.Get("relay")
			if relay == nil {
				continue
			}

			v, err := relay.Get()
			if err != nil {
				otto.GetLogger().Error("Error getting input value: ", "error", err.Error())
				continue
			}

			if v == 0 {
				v = 1
			} else {
				v = 0
			}
			fmt.Printf("v: %+v\n", v)
			val := strconv.Itoa(v)
			for _, t := range dev.Pubs {
				otto.GetMQTT().Publish(t, val)
			}

		case <-done:
			return
		}
	}
}

func (dm *DeviceManager) initDevices(done chan bool) {
	led := GPIOOut("led", 6)
	relay := GPIOOut("relay", 22)
	button := GPIOIn("button", 23)
	devices.Add(led)
	devices.Add(relay)
	devices.Add(button)

	// Subscribe to the button if it is pushed we want to know
	// about it and turn on and off the LED and Relay simultaneously
	for _, pub := range button.Pubs {
		m := otto.GetMQTT()
		m.Subscribe(pub, dm)
	}

	go ButtonHandler()

}
