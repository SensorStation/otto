package main

import (
	"embed"
	"flag"

	"github.com/sensorstation/otto"
	"github.com/sensorstation/otto/cmd"
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/devices/bme280"
	"github.com/sensorstation/otto/devices/button"
	"github.com/sensorstation/otto/devices/led"
	"github.com/sensorstation/otto/devices/relay"
	"github.com/sensorstation/otto/message"
	"github.com/sensorstation/otto/mocks"
)

var (
	l        *otto.Logger
	done     chan bool
	mock     bool
	mockMQTT bool
	mockGPIO bool

	cli bool
)

func init() {
	flag.BoolVar(&mockMQTT, "mock", false, "mock the hardware")
	flag.BoolVar(&mockMQTT, "mock-mqtt", false, "mock the hardware")
	flag.BoolVar(&mockGPIO, "mock-gpio", false, "mock the hardware")
	flag.BoolVar(&cli, "cli", false, "Run the otto interactive cli")
}

func main() {
	flag.Parse()

	l = otto.GetLogger()
	done = make(chan bool)

	if mock {
		mockMQTT = true
		mockGPIO = true
	}
	if mockMQTT {
		otto.GetMQTTClient(mocks.GetMockClient())
	}
	if mockGPIO {
		devices.GetGPIO().Mock = true
	}

	// TODO capture signals
	initSignals()
	initApp()
	initStations()
	initDevices(done)
	if cli {
		cmd.Execute()
	}

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
	go butOn.EventLoop(done)
	m.SubscribeHandle(otto.TopicControl("on"), func(msg *message.Msg) {
		m.Publish(otto.TopicControl("relay"), "on")
		m.Publish(otto.TopicControl("led"), "on")
	})

	butOff := button.New("off", 27)
	go butOff.EventLoop(done)
	m.SubscribeHandle(otto.TopicControl("off"), func(msg *message.Msg) {
		m.Publish(otto.TopicControl("relay"), "off")
		m.Publish(otto.TopicControl("led"), "off")
	})

	bme := bme280.New("bme", "/dev/i2c-1", 0x76)
	if bme == nil {
		return
	}
	bme.Pubs = append(bme.Pubs, otto.TopicData("bme280"))
	bme.Period = 10
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
