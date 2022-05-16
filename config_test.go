package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig(t *testing.T) {
	flag.Parse()

	fname := "/tmp/sensor-test.json"

	// config will already been configured
	err := config.SaveFile(fname)
	if err != nil {
		t.Error("saving config: ", err)
	}

	var c Configuration
	err = c.ReadFile(fname)
	if err != nil {
		t.Error("reading config: ", err)
	}

	if c != config {
		t.Errorf("c (%+v) != config (%+v)", c, config)
	}
}

func TestConfigHTTP(t *testing.T) {

	ts := httptest.NewServer(Configuration{})
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	cbuf, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	var c Configuration
	err = json.Unmarshal(cbuf, &c)
	if err != nil {
		t.Error(err)
	}

	if c != config {
		t.Errorf("c (%+v) != config (%+v)", c, config)
	}
}
