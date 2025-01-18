package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sensorstation/otto/devices/oled"
	"periph.io/x/devices/v3/ssd1306"
)

// const On = ssd1306.On
// const Off = ssd1306.Off

var (
	display *oled.OLED
	err     error
)

func main() {

	display, err = oled.New("oled", 128, 64)
	if err != nil {
		log.Fatal(err)
	}

	// println(display.Dev.String())
	foo()

	// var stime time.Duration = time.Second * 2
	// for {
	// 	hello()
	// 	time.Sleep(stime)

	// ballerine()

	// 	sensors()
	// 	time.Sleep(stime)
	// }

}

func foo() {
	display.Clear()
	display.DrawString(10, 10, "whoops")
	display.Diagonal(20, 10, 100, 20, oled.On)
	// display.Diagonal(60, 20, 60, 60, oled.On)

	// display.Diagonal(10, 30, 110, 50, oled.On)
	// display.Diagonal(10, 30, 10, 40, oled.On)

	// display.Diagonal(10, 40, 110, 60, oled.On)
	// display.Diagonal(110, 50, 110, 60, oled.On)

	display.Draw()
}

func hello() {
	// draw a lgine at 50, 100 lenght 50 pixels, 4 wide
	display.Clear()
	display.Line(0, 12, display.Width, 2, oled.On)
	display.Rectangle(100, 40, 120, 60, oled.On)
	display.DrawString(10, 10, "Hello, world!")
	display.Draw()
}

func ballerine() {
	var done <-chan time.Time
	done = time.After(10 * time.Second)

	display.Clear()
	display.AnimatedGIF("ballerine.gif", done)
	display.Draw()
}

func sensors() {

	display.Clear()

	temp := 10.1
	pressure := 11.2
	humidity := 12.3
	//relay := "On"

	start := 25
	t := time.Now().Format(time.Kitchen)
	display.DrawString(10, 10, "OttO: "+t)
	display.DrawString(10, start, fmt.Sprintf("temp: %7.2f", temp))
	display.DrawString(10, start+15, fmt.Sprintf("pres: %7.2f", pressure))
	display.DrawString(10, start+30, fmt.Sprintf("humi: %7.2f", humidity))

	display.Draw()
}

func scroller(done chan bool) {
	sensors()
	display.Scroll(ssd1306.Right, ssd1306.FrameRate5, 0, -1)

	<-done
	display.StopScroll()

}
