package main

import (
	"embed"
	"flag"

	"github.com/sensorstation/otto/cmd"
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/messanger"
	"github.com/sensorstation/otto/server"
	"github.com/sensorstation/otto/utils"
)

var (
	done     chan bool
	loglevel string
	mock     bool
	mockMQTT bool
	mockGPIO bool

	cli bool
)

func init() {
	flag.BoolVar(&cli, "cli", false, "Run the otto interactive cli")
	flag.BoolVar(&mockMQTT, "mock", false, "mock the hardware")
	flag.BoolVar(&mockMQTT, "mock-mqtt", false, "mock the hardware")
	flag.BoolVar(&mockGPIO, "mock-gpio", false, "mock the hardware")
	flag.StringVar(&loglevel, "log", "warn", "Logging level debug, info, warn, error, fatal")
}

func main() {
	flag.Parse()

	done = make(chan bool)
	initLogging()
	initMock()

	// TODO capture signals
	c := &controller{}
	c.initApp()
	c.initSignals()
	c.initStations()
	c.initDataManager()
	c.initDevices(done)

	if cli {
		cmd.Execute()
	}

	<-done
	c.cleanup()
}

func initLogging() {

	level := loglevel
	logfile := "otto.log"
	utils.InitLogger(level, logfile)
}

func initMock() {
	if mock {
		mockMQTT = true
		mockGPIO = true
	}
	if mockMQTT {
		messanger.SetMQTTClient(messanger.GetMockClient())
	}
	if mockGPIO {
		devices.GetGPIO().Mock = true
	}
}

//go:embed app
var content embed.FS

func (c *controller) initApp() {
	s := server.GetServer()

	// The following line is commented out because
	var data any
	s.EmbedTempl("/emb", data, content)
	s.Appdir("/", "app")
	go s.Start()
	c.Server = s
}
