package otto

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Configuration holds all variables that can be changed
// to alter the way the program works.
type Configuration struct {
	Addr   string
	Broker string
	Appdir string

	Debug       bool
	DebugMQTT   bool
	FakeWS      bool
	Filename    string
	Interactive bool
	Mock        bool
	MaxData     int // maximum data values to save
	Plugin      string

	GPIO    bool
	Verbose bool
}

// ServeHTTP provides a REST interface to the config structure
func (c Configuration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(config)

	case "POST", "PUT":
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
