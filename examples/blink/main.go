/*
Blink sets up pin 6 for an LED and goes into an endless
toggle mode.
*/

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sensorstation/otto/device"
	"github.com/sensorstation/otto/devices/led"
	"github.com/sensorstation/otto/messanger"
)

var (
	useMQTT bool
	pinid   int
	mock    string
	count   int
)

func init() {
	flag.BoolVar(&useMQTT, "mqtt", false, "Use mqtt or a timer")
	flag.IntVar(&pinid, "pin", 6, "The GPIO pin the LED is attached to")
	flag.StringVar(&mock, "mock", "", "mock gpio and/or mqtt")
}

func main() {
	flag.Parse()

	// Create the led, name it "led" and add a publish topic
	led, done := initLED("led", pinid)
	if useMQTT {
		domqtt(led)
	} else {
		dotimer(led, 1*time.Second, done)
		fmt.Println("LED will blink every second")
	}
	<-done
}

func initLED(name string, pin int) (*led.LED, chan any) {
	led := led.New(name, pin)
	led.AddPub(messanger.TopicData(led.Name()))
	done := make(chan any)
	return led, done
}

func domqtt(led *led.LED) {
	led.Subscribe(messanger.TopicControl(led.Name()), led.Callback)
}

func dotimer(led *led.LED, period time.Duration, done chan any) {
	count = 0
	led.TimerLoop(period, done, func() error {
		// led.Set(count % 2)
		count++
		return nil
	})
}

func domock() {
	switch mock {
	case "mqtt":
		messanger.SetMQTTClient(messanger.GetMockClient())

	case "gpio":
		device.GetGPIO().Mock = true

	case "both":
		messanger.SetMQTTClient(messanger.GetMockClient())
		device.GetGPIO().Mock = true

	default:
		return
	}
}
