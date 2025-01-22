package main

import (
	"fmt"
	"log"

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
