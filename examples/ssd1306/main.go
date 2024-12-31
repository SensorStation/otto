package main

import (
	"log"
	"time"

	"github.com/sensorstation/otto/devices/ssd1306"
)

func main() {

	display, err := ssd1306.NewDisplay("oled", 128, 64)
	if err != nil {
		log.Fatal(err)
	}
	// draw a line at 50, 100 lenght 50 pixels, 4 wide
	display.Clear()
	display.Line(0, 12, display.Width, 2, ssd1306.On)
	display.Rectangle(100, 40, 120, 60, ssd1306.On)
	display.DrawString(10, 10, "Hello, world!")
	display.Draw()

	time.Sleep(10 * time.Second)

	display.Clear()
	display.AnimatedGIF("ballerine.gif")
	display.Draw()
}
