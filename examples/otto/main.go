package main

import (
	"embed"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/devices/bme280"
	"github.com/sensorstation/otto/devices/button"
	"github.com/sensorstation/otto/devices/led"
	"github.com/sensorstation/otto/devices/relay"
)

var (
	l    *otto.Logger
	done chan bool
)

func main() {
	l = otto.GetLogger()
	done = make(chan bool)

	// TODO capture signals

	initSignals()
	initApp()
	initStations()
	initDevices(done)
	// cmd.Execute()

	<-done
	cleanup()
}

func cleanup() {
	g := devices.GetGPIO()
	g.Shutdown()
}

func initSignals() {
	// todo make sure we capture signals
}

func initStations() {
	st := otto.GetStationManager()
	st.Start()
}

func initDevices(done chan bool) {

	m := otto.GetMQTT()
	relay := relay.New("relay", 22)
	m.Subscribe(otto.TopicControl("relay"), relay)

	led := led.New("led", 6)
	m.Subscribe(otto.TopicControl("led"), led)

	butOn := button.New("on", 23)
	butOn.Pubs = append(butOn.Pubs, otto.TopicControl("on"))
	go butOn.EventLoop(done)

	butOff := button.New("off", 27)
	go butOff.EventLoop(done)

	bme := bme280.New("bme", "/dev/i2c-1", 0x76)
	err := bme.Init()
	if err != nil {
		otto.GetLogger().Error("Failed to initialize bme", "error", err)
		return
	}

	bme.Pubs = append(bme.Pubs, otto.TopicData("bme280"))
	go bme.Loop(done)

}

//go:embed app
var content embed.FS

func initApp() {
	s := otto.GetServer()

	// The following line is commented out because
	var data any
	s.EmbedTempl("/emb", data, content)
	s.Appdir("/", "app")
	go s.Start()
}
