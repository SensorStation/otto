package gtu7

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/messanger"
)

type GTU7 struct {
	*devices.SerialDevice
	lastGPS GPS

	scanner *bufio.Scanner
}

func NewGTU7(devname string) *GTU7 {
	g := &GTU7{}
	g.SerialDevice = devices.NewSerialDevice("gt-u7", devname, 9600)
	g.AddPub(messanger.TopicData("gt-u7"))
	return g
}

func (g *GTU7) Open() error {

	err := g.SerialDevice.Open()
	if err != nil {
		fmt.Printf("Failed to open serial port %s - %v\n", g.PortName, err)
		return err
	}
	return nil
}

func (g *GTU7) OpenRead() error {
	err := g.Open()
	if err != nil {
		return err
	}
	g.scanner = bufio.NewScanner(g.Port)
	return nil
}

func (g *GTU7) OpenStrings(input string) {
	g.scanner = bufio.NewScanner(strings.NewReader(input))
}

func (g *GTU7) StartReading() chan *GPS {
	l := logger.GetLogger()

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
			l.Error("scanning GPS data", "error", err)
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
