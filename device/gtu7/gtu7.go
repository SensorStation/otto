package gtu7

import (
	"bufio"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sensorstation/otto/device"
)

type GTU7 struct {
	*device.Device
	*device.Serial
	lastGPS GPS
	scanner *bufio.Scanner
}

func NewGTU7(devname string) *GTU7 {
	g := &GTU7{}
	g.Device = device.NewSerialDevice("gt-u7", devname, 9600)
	return g
}

func (g *GTU7) Open() error {
	err := g.Device.Open()
	if err != nil {
		fmt.Printf("Failed to open serial port %s - %v\n", g.String(), err)
		return err
	}
	return nil
}

func (g *GTU7) OpenRead() error {
	err := g.Open()
	if err != nil {
		return err
	}
	g.scanner = bufio.NewScanner(g)
	return nil
}

func (g *GTU7) OpenStrings(input string) {
	g.scanner = bufio.NewScanner(strings.NewReader(input))
}

func (g *GTU7) StartReading() chan *GPS {
	parseQ := make(chan string)
	gpsQ := g.startParser(parseQ)
	go func() {
		for g.scanner.Scan() {
			if g.scanner.Text() == "" {
				continue
			}
			line := g.scanner.Text()
			parseQ <- line
		}
		if err := g.scanner.Err(); err != nil {
			slog.Error("scanning GPS data", "error", err)
		}
		close(parseQ)
	}()
	return gpsQ
}

func (g *GTU7) startParser(parseQ chan string) chan *GPS {

	gps := &GPS{}
	gpsQ := make(chan *GPS)
	go func() {
		for line := range parseQ {
			gps.Parse(line)
			if gps.IsComplete() {
				gpsQ <- gps
			}
		}
		close(gpsQ)
	}()
	return gpsQ
}
