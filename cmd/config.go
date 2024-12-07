package cmd

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
)

var (
	config Configuration
)

// Configuration holds all variables that can be changed
// to alter the way the program works.
type Configuration struct {
	Addr   string
	Broker string
	Appdir string

	Debug     bool
	DebugMQTT bool
	FakeWS    bool
	Filename  string
	Mock      bool
	MaxData   int // maximum data values to save
	Plugin    string

	GPIO    bool
	Verbose bool
}

func init() {
	flag.StringVar(&config.Addr, "addr", "0.0.0.0:8011", "Address to listen for web connections")
	flag.StringVar(&config.Appdir, "appdir", "/srv/otto/dist", "Directory for the web app distribution")
	flag.StringVar(&config.Broker, "broker", "localhost", "Address of MQTT broker")
	flag.BoolVar(&config.Debug, "debug", false, "Start debugging")
	flag.BoolVar(&config.DebugMQTT, "debug-mqtt", false, "Debugging MQTT messages")
	flag.BoolVar(&config.FakeWS, "fake-ws", false, "Fake websocket data")
	flag.StringVar(&config.Filename, "config", "~/.config/sensors.json", "Where to read and store config")
	flag.IntVar(&config.MaxData, "max-data", 1000, "Maximum data length for sensors")
	flag.BoolVar(&config.Mock, "mock", false, "Mock sensor data")
	flag.StringVar(&config.Plugin, "plugin", "", "Plugins to be loaded")
	flag.BoolVar(&config.Verbose, "verbose", false, "Crank up the output")
	flag.BoolVar(&config.GPIO, "gpio", false, "Utilize GPIO for Raspberry PI")
}

// func GetConfig() Configuration {
// 	return config
// }

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
		l.Error("JSON marshaling config", "error", err)
		return err
	}

	err = ioutil.WriteFile(fname, jbuf, 0644)
	if err != nil {
		l.Error("FILE writing config", "error", err)
		return err
	}
	return err
}

// Load the file from the file corresponding to the fname parameter
func (c *Configuration) ReadFile(fname string) error {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		l.Error("failed to read", "file", fname, "error", err)
		return err
	}

	err = json.Unmarshal(buf, c)
	if err != nil {
		l.Error("failed to read ", "file", fname, "error", err)
		return err
	}

	return err
}

func (c Configuration) GetAddr() string {
	return c.Addr
}

func (c Configuration) GetAppdir() string {
	return c.Appdir
}

func (c Configuration) GetBroker() string {
	return c.Broker
}
