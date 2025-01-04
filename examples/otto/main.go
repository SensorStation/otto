package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/sensorstation/otto/cmd"
	"github.com/sensorstation/otto/data"
	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/devices/bme280"
	"github.com/sensorstation/otto/devices/button"
	"github.com/sensorstation/otto/devices/led"
	"github.com/sensorstation/otto/devices/relay"
	"github.com/sensorstation/otto/devices/ssd1306"
	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/message"
	"github.com/sensorstation/otto/messanger"
	"github.com/sensorstation/otto/server"
	"github.com/sensorstation/otto/station"
)

var (
	l        *logger.Logger
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

	l = logger.GetLogger()
	done = make(chan bool)

	if mock {
		mockMQTT = true
		mockGPIO = true
	}
	if mockMQTT {
		messanger.GetMQTTClient(messanger.GetMockClient())
	}
	if mockGPIO {
		devices.GetGPIO().Mock = true
	}

	// TODO capture signals
	initSignals()
	initApp()
	initStations()
	initDataManager()

	err := initDevices(done)
	if err != nil {
		fmt.Println("Failed to initialize devices: ", err)
	}

	if cli {
		cmd.Execute()
	}

	<-done
	cleanup()
}

type controller struct {
}

func (c *controller) init() {

}

func cleanup() {
	g := devices.GetGPIO()
	g.Shutdown()
}

func initSignals() {
	// todo make sure we capture signals
}

func initDataManager() {
	dm := data.GetDataManager()
	dm.Subscribe("ss/d/#")

	srv := server.GetServer()
	srv.Register("/api/data", dm)
}

func initStations() {
	st := station.GetStationManager()
	st.Start()
}

func initDevices(done chan bool) error {
	initRelay(22)

	initLED(6)
	initButton("on", 23)
	initButton("off", 27)
	bme, err := initBME280("/dev/i2c-1", 0x76, done)
	if err != nil {
		return err
	}
	bme.AddPub(messanger.TopicData("bme280"))
	bme.Period = 10 * time.Second
	go bme.TimerLoop(done, bme.ReadPub)

	initOLED(done)
	return err
}

func initRelay(idx int) {
	relay := relay.New("relay", idx)
	relay.AddPub(messanger.TopicData("relay"))
	relay.Subscribe(messanger.TopicControl("button"), relay.Callback)
}

func initLED(idx int) {
	led := led.New("led", idx)
	led.AddPub(messanger.TopicData("led"))
	led.Subscribe(messanger.TopicControl("button"), led.Callback)
}

func initButton(name string, idx int) {
	but := button.New(name, idx)
	but.AddPub(messanger.TopicControl("button"))
	go but.EventLoop(done, but.ReadPub)
}

func initBME280(bus string, addr int, done chan bool) (bme *bme280.BME280, err error) {
	bme = bme280.New("bme280", "/dev/i2c-1", 0x76)
	if bme == nil {
		return nil, fmt.Errorf("Failed initialize BME280 %s %d", "/dev/i2c-1", 0x76)
	}

	if mockGPIO {
		bme.Mock = true
	}
	err = bme.Init()
	if err != nil {
		return nil, err
	}
	return bme, nil
}

func initOLED(done chan bool) {

	display, err := ssd1306.NewDisplay("oled", 128, 64)
	if err != nil {
		log.Fatal(err)
	}

	dispvals := struct {
		temp     float64
		pressure float64
		humidity float64
		relay    string
	}{}

	refresh := func() {
		display.Clear()
		start := 25
		t := time.Now().Format(time.Kitchen)
		display.DrawString(10, 10, "OttO: "+t)
		display.DrawString(10, start, fmt.Sprintf("temp: %7.2f", dispvals.temp))
		display.DrawString(10, start+15, fmt.Sprintf("pres: %7.2f", dispvals.pressure))
		display.DrawString(10, start+30, fmt.Sprintf("humi: %7.2f", dispvals.humidity))
		display.Draw()

	}

	m := messanger.GetMQTT()
	m.Subscribe(messanger.TopicData("bme280"), func(msg *message.Msg) {
		mm, err := msg.Map()
		if err != nil {
			l.Error("Failed top get map", "error", err)
		}

		var ex bool
		dispvals.temp, ex = mm["Temperature"].(float64)
		if !ex {
			l.Error("failed to get temperature")
		}
		dispvals.pressure, ex = mm["Pressure"].(float64)
		if !ex {
			l.Error("failed to get pressure")
		}
		dispvals.humidity, ex = mm["Humidity"].(float64)
		if !ex {
			l.Error("failed to get Humidity")
		}
		refresh()
	})

	m.Subscribe(messanger.TopicData("relay"), func(msg *message.Msg) {
		dispvals.relay = msg.String()
	})
}

//go:embed app
var content embed.FS

func initApp() {
	s := server.GetServer()

	// The following line is commented out because
	var data any
	s.EmbedTempl("/emb", data, content)
	s.Appdir("/", "app")
	go s.Start()
}
