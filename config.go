package main

import (
	"flag"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

type Config interface {
	http.Handler
	SaveFile(fname string) error
	ReadFile(fname string) error
}

type Configuration struct {
	Addr   string
	App    string
	Broker string

	Debug		bool
	DebugMQTT	bool
	FakeWS		bool
	Filename	string
	Mock		bool

	MaxData	int					// maximum data values to save

	GPIO	bool
	Verbose bool
}

var (
	config	 Configuration
)

func init() {
	flag.StringVar(&config.Addr, "addr", "0.0.0.0:8011", "Address to listen for web connections")
	flag.StringVar(&config.App, "app", "../app/dist", "Directory for the web app distribution")
	flag.StringVar(&config.Broker, "broker", "localhost", "Address of MQTT broker")
	flag.BoolVar(&config.Debug, "debug", false, "Start debugging")
	flag.BoolVar(&config.DebugMQTT, "debug-mqtt", false, "Debugging MQTT messages")
	flag.BoolVar(&config.FakeWS, "fake-ws", false, "Fake websocket data")
	flag.IntVar(&config.MaxData, "max-data", 1000, "Maximum data length for sensors")
	flag.BoolVar(&config.Mock, "mock", false, "Mock sensor data")
	flag.StringVar(&config.Filename, "config", "~/.config/sensors.json", "Where to read and store config")

	flag.BoolVar(&config.Verbose, "verbose", false, "Crank up the output")
	flag.BoolVar(&config.GPIO, "gpio", false, "Utilize GPIO for Raspberry PI")
}

func GetConfig() Configuration {
	return config
}

// ServeHTTP provides a REST interface to the config structure
func (c Configuration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(config)

	case "POST", "PUT":
		// TODO
		http.Error(w, "Not Yet Supported", 401)
	}
}

// Save write the configuration to a file in JSON format
func (c *Configuration) SaveFile(fname string) error {

	jbuf, err := json.Marshal(c)
	if err != nil {
		log.Printf("[ERROR]: JSON marshaling config: %+v", err)
		return err
	}

	err = ioutil.WriteFile(fname, jbuf, 0644)
	if err != nil {
		log.Printf("[ERROR]: FILE writing config: %+v", err)
		return err
	}
	return err
}

// Load the file from the file corresponding to the fname parameter
func (c *Configuration) ReadFile(fname string) error {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Printf("[ERROR]: failed to read file %s, %v", fname, err)
		return err
	}

	err = json.Unmarshal(buf, c)
	if err != nil {
		log.Printf("[ERROR]: failed to read file %s, %v", fname, err)
		return err
	}

	return err
}
