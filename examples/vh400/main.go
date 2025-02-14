package main

import (
	"fmt"

	"github.com/sensorstation/otto/device/vh400"
)

func main() {
	var readQ [4]<-chan float64
	for i := 0; i < 4; i++ {
		s := vh400.New("vh400", i)
		readQ[i] = s.AnalogPin.ReadContinuous()
	}

	for {
		var val any
		var idx int
		select {
		case val = <-readQ[0]:
			idx = 0
		case val = <-readQ[1]:
			idx = 1
		case val = <-readQ[2]:
			idx = 2
		case val = <-readQ[3]:
			idx = 3
		}
		fmt.Printf("%d: %5.2f\n", idx, val.(float64))
	}
}
