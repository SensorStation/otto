package main

import (
	"fmt"

	"github.com/sensorstation/otto/devices/vh400"
)

func main() {
	var readQ [4]<-chan float64
	for i := 0; i < 4; i++ {
		s := vh400.New("vh400", i)
		readQ[i] = s.ReadContinuous()
	}

	for {
		select {
		case val := <-readQ[0]:
			fmt.Printf("0: %5.2f\n", val)
		case val := <-readQ[1]:
			fmt.Printf("1: %5.2f\n", val)
		case val := <-readQ[2]:
			fmt.Printf("2: %5.2f\n", val)
		case val := <-readQ[3]:
			fmt.Printf("3: %5.2f\n", val)
		}
	}
}
