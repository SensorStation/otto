package iote

import (
	"log"

	"encoding/json"
	"io/ioutil"
	"net/http"

)

type Configuration struct {
	Addr		string		`json:"addr"`
	Broker		string		`json:"mqtt-broker"`
	DebugMQTT	bool		`json:"mqtt-debug"`
	MaxData		int			`json:"max-data"`
	Verbose		bool		`json:"verbose"`
}

var (
	config Configuration
)

func init() {
	config.Addr		= ":8011"
	config.Broker	= "tcp://localhost:1883"
}

func GetConfig() (Configuration) {
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
