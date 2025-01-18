package main

import (
	"fmt"
	"log"

	"github.com/sensorstation/otto/data"
	"github.com/sensorstation/otto/devices/ads1115"
	"github.com/sensorstation/otto/messanger"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

const (
	adcRange = 26400 // ADS1115 gives you 32768 steps with 4.096v, with 3.3v which is the 80% of 4.096v we get 26400 steps
	// channel               = ads1x15.Channel3
	channel               = ads1x15.Channel0
	inputVoltageInChannel = 3300 * physic.MilliVolt
)

func main() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open default I²C bus.
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer bus.Close()
	fmt.Printf("bus: %+v\n", bus)

	// Create a new ADS1115 ADC.
	adc, err := ads1x15.NewADS1115(bus, &ads1x15.DefaultOpts)
	if err != nil {
		log.Fatalln(err)
	}

	// Obtain an analog pin from the ADC.
	pin, err := adc.PinForChannel(channel, inputVoltageInChannel, 1*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		log.Fatalln(err)
	}
	defer pin.Halt()

	fmt.Printf("%+v\n", pin)
	reading, err := pin.Read()
	fmt.Printf("reading: %+v: %+v\n", reading, err)

	// Read values continuously from ADC.
	fmt.Println("Continuous reading")
	c := pin.ReadContinuous()

	for reading := range c {
		// voltage := float32(reading.Raw) / (adcRange) * 3.3
		// Rt := 10 * voltage / (3.3 - voltage)
		// tempK := 1 / (1/(273.15+25) + math.Log(float64(Rt/10))/3950.0) // calculate temperature in Kelvin
		// tempC := tempK - 273.15                                        // calculate temperature (Celsius)
		// fmt.Printf("Temp in Kelvin: %0.2f\n", tempK)
		// fmt.Printf("Temp in Celsius: %0.2f\n", tempC)
		fmt.Println(reading)
	}
}

func omain() {
	// create the topic the bme will publish to and the DataManager
	// will subscribe to
	topic := messanger.TopicData("ads1115")

	// Set the BME i2c device and address Initialize the bme to use
	// the i2c bus
	ads := ads1115.New("ads1115", "/dev/i2c-1", 0x48)
	ads.AddPub(topic)
	err := ads.Init()
	if err != nil {
		panic(err)
	}

	err = ads.Read()
	if err != nil {
		panic(err)
	}

	// Before we start reading temp, etc. let's subscribe to
	// the messages we are going to publish.
	dm := data.GetDataManager()
	dm.Subscribe(topic, dm.Callback)

	// start reading in a loop and publish the results via MQTT
	done := make(chan bool)
	go ads.TimerLoop(done, ads.ReadPub)
	<-done
}
