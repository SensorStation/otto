package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sensorstation/otto/data"
	"github.com/sensorstation/otto/logger"
	"github.com/sensorstation/otto/messanger"
)

var (
	path = "ss/d/station/env3"
)

func main() {
	// Set the BME i2c device and address
	bme := BME280{
		Addr: 0x76,
		Dev:  "/dev/i2c-1",
	}

	// Initialize the bme to use the i2c bus
	err := bme.Init()
	if err != nil {
		panic(err)
	}

	// Get mqtt ready to start publishing the results
	// from the bme via mqtt
	mqtt := messanger.GetMQTT()

	// Before we start reading temp, etc. let's subscribe to
	// the messages we are going to publish.
	dm := data.GetDataManager()
	mqtt.Subscribe(path, dm.Callback)

	// start reading in a loop and publish the results
	// via MQTT
	done := make(chan bool)
	go func() {
		for {
			// Read temp, humidity and pressure
			vals, err := bme.Read()
			if err != nil {
				panic(err)
			}

			// Print the values
			fmt.Printf("vals: %+v\n", vals)
			jb, err := json.Marshal(vals)
			if err != nil {
				logger.GetLogger().Error("failed to unmarshal bme Response", "error", err.Error())
				done <- true
				break
			}

			mqtt.Publish(path, jb)
			time.Sleep(1 * time.Second)
		}
	}()
	<-done
}

type BMEReader struct {
	Path string
}
