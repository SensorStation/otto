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

func cleanup() {
	g := devices.GetGPIO()
	g.Shutdown()
}

func initDataManager() {
	dm := data.GetDataManager()
	srv := server.GetServer()
	srv.Register("/api/data", dm)
}

func initSignals() {
	// todo make sure we capture signals
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
	initOLED(bme, done)
	return err
}

func initRelay(idx int) {
	m := messanger.GetMQTT()
	relay := relay.New("relay", idx)
	m.Subscribe(messanger.TopicControl("relay"), relay.Callback)
}

func initLED(idx int) {
	m := messanger.GetMQTT()
	led := led.New("led", idx)
	m.Subscribe(messanger.TopicControl("led"), led.Callback)
}

func initButton(name string, idx int) {
	m := messanger.GetMQTT()
	but := button.New(name, idx)
	go but.EventLoop(done)
	m.Subscribe(messanger.TopicControl(name), func(msg *message.Msg) {
		m.Publish(messanger.TopicControl("relay"), name)
		m.Publish(messanger.TopicControl("led"), name)
	})
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
	bme.Period = 10 * time.Second
	go bme.Loop(done)
	return bme, nil
}

func initOLED(bme *bme280.BME280, done chan bool) {

	display, err := ssd1306.NewDisplay("oled", 128, 64)
	if err != nil {
		log.Fatal(err)
	}

	m := messanger.GetMQTT()
	m.Subscribe(bme.Pub, func(msg *message.Msg) {
		mm, err := msg.Map()
		if err != nil {
			l.Error("Failed top get map", "error", err)
		}

		temp, ex := mm["Temperature"]
		if !ex {
			l.Error("failed to get temperature")
		}
		pressure, ex := mm["Pressure"]
		if !ex {
			l.Error("failed to get pressure")
		}
		humidity, ex := mm["Humidity"]
		if !ex {
			l.Error("failed to get Humidity")
		}

		display.Clear()
		start := 25
		display.DrawString(10, 10, "OttO IoT Champion")
		display.DrawString(10, start, fmt.Sprintf("temp: %7.2f", temp))
		display.DrawString(10, start + 15, fmt.Sprintf("pres: %7.2f", pressure))
		display.DrawString(10, start + 30, fmt.Sprintf("humi: %7.2f", humidity))
		display.Draw()
	})

	// go func() {
	// 	running := true
	// 	for running {
	// 		select {
	// 		case <-done:
	// 			running = false

	// 		default:
	// 			// draw a line at 50, 100 lenght 50 pixels, 4 wide
	// 			display.Clear()
	// 			display.Line(0, 12, display.Width, 2, ssd1306.On)
	// 			display.Rectangle(100, 40, 120, 60, ssd1306.On)
	// 			display.DrawString(10, 10, "Hello, world!")
	// 			display.Draw()

	// 			time.Sleep(10 * time.Second)

	// 			display.Clear()
	// 			display.AnimatedGIF("ballerine.gif")
	// 			display.Draw()
	// 			time.Sleep(10 * time.Second)
	// 		}
	// 	}
	// }()
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
