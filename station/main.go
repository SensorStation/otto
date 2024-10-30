package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sensorstation/otto/gpio"
)

// GPIO16 DHT22
// GPIO 6 LED

func main() {

	gpio := gpio.GetRPI()
	led := gpio.PinInit("green-led", 6, gpio.Output(0))
	// dht := rpi.Pin(16, "am2302", rpi.Input)

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

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

/*
func main0() {
	offset := 6
	chip := "gpiochip4"
	v := 0
	l, err := gpiocdev.RequestLine(chip, offset, gpiocdev.AsOutput(v))
	if err != nil {
		panic(err)
	}
	// revert line to input on the way out.
	defer func() {
		l.Reconfigure(gpiocdev.AsInput)
		fmt.Printf("Input pin %s:%d\n", chip, offset)
		l.Close()
	}()
	values := map[int]string{0: "inactive", 1: "active"}
	fmt.Printf("Set pin %s:%d %s\n", chip, offset, values[v])

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	for {
		select {
		case <-time.After(500 * time.Millisecond):
			v ^= 1
			l.SetValue(v)
			fmt.Printf("Set pin %s:%d %s\n", chip, offset, values[v])
		case <-quit:
			return
		}
	}
}
*/
