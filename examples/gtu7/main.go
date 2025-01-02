package main

import "github.com/sensorstation/otto/devices/gtu7"

func main() {
	g := gtu7.NewGTU7("/dev/ttyS0")
	gpsQ := g.StartReading()

	for gps := range gpsQ {
		g.Publish(gps)
	}
}
