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
	stationName := otto.StationName
	dm := devices.GetDeviceManager()

	relay := relay.New("relay", 22)
	m.Subscribe("ss/c/"+stationName+"/off", relay)
	m.Subscribe("ss/c/"+stationName+"/on", relay)
	dm.Add(relay)

	led := led.New("led", 6)
	m.Subscribe("ss/c/"+stationName+"/off", led)
	m.Subscribe("ss/c/"+stationName+"/on", led)
	dm.Add(led)

	onButton := button.New("on", 23)
	dm.Add(onButton)

	offButton := button.New("off", 27)
	dm.Add(offButton)
	go onButton.EventLoop(done)
	go offButton.EventLoop(done)

	bme := bme280.New("bme", "/dev/i2c-1", 0x76)
	err := bme.Init()
	if err != nil {
		l.Error("Failed to open the bme280", "error", err)
		return
	}
	dm.Add(bme)
	bme.AddPub("ss/d/" + stationName + "/bme280")
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
