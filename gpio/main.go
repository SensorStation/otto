package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sensorstation/otto"
)

var (
	config *Configuration
	e      Echo
	mqtt   *otto.MQTT
)

// DHT22		- GPIO 16 DHT22
// Green LED	- GPIO 6
func init() {
	config = &Configuration{}
	flag.StringVar(&config.Broker, "broker", "localhost", "MQTT Broker")
}

func main() {
	flag.Parse()

	mqtt := &otto.MQTT{
		Broker: config.Broker,
	}
	mqtt.Start()
	mqtt.Subscribe("ss/echo", e)
	mqtt.Subscribe("ss/echo2", e2)

	run()
}

type Echo struct {
}

func (e Echo) Callback(t string, payload []byte) {
	log.Println("echo: ", t, string(payload))
}

func run() {

	// Open GPIO
	gpio := GetGPIO()

	// Configure all the devices
	led := gpio.Pin("green-led", 6, Output(0))
	// dht := rpi.PinInit("input", 16, ModeInput)

	// revert line to input on the way out.
	defer func() {
		gpio.Shutdown()
	}()

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	// fmt.Println(gpio.String())

	// Start the loop
	v := 0
	for {
		select {
		case <-time.After(1 * time.Second):
			v ^= 1
			led.Set(v)
			fmt.Printf("LED %s\n", led.String())

		case <-quit:
			return
		}

	}
}
