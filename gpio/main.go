package gpio

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
	mqtt   *otto.MQTT
	gpio   *GPIO

	e Echo
)

func main() {
	mqtt = &otto.MQTT{
		Broker: config.Broker,
	}
	mqtt.Start()
	mqtt.Subscribe("ss/echo", e)

	run()
}

func run() {

	// Open GPIO
	gpio = GetGPIO()

	// Configure all the devices
	led := gpio.Pin("green-led", 6, Output(0))
	dht := NewDHT("am2302", 16, ModeInput)

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
		case <-time.After(5 * time.Second):
			v ^= 1
			led.Set(v)
			fmt.Printf("LED %s\n", led.String())

			err := dht.Read()
			if err != nil {
				fmt.Printf("dht err: %v\n", err)
			}
			fmt.Printf("temperature: %5.2f - humidity: %5.2f\n",
				dht.Temperature(), dht.Humidity())

		case <-quit:
			return
		}

	}
}

type Echo struct {
}

func (e Echo) Callback(t string, payload []byte) {
	log.Println("echo: ", t, string(payload))
}
