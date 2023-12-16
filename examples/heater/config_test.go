package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkDefaults(t *testing.T, c *Configuration) bool {
	if c.Addr != "0.0.0.0:8011" {
		t.Errorf("config.Addr expected (0.0.0.0:8011) got (%s)", c.Addr)
		return false
	}

	if c.App != "../app/dist" {
		t.Errorf("config.App expected (../app/dist) got (%s)", c.App)
		return false
	}

	if c.Broker != "localhost" {
		t.Errorf("config.Broker expected (localhost) got (%s)", c.Broker)
		return false
	}

	if c.Debug != false {
		t.Errorf("config.Debug expected (false) got (%t)", c.Debug)
		return false
	}

	if c.DebugMQTT != false {
		t.Errorf("config.DebugMQTT expected (false) got (%t)", c.DebugMQTT)
		return false
	}

	if c.FakeWS != false {
		t.Errorf("config.FakeWS expected (false) got (%t)", c.FakeWS)
		return false
	}

	if c.Filename != "~/.config/sensors.json" {
		t.Errorf("config.Filename expected (~/.config/sensors.json) got (%s)", c.Filename)
		return false
	}

	if c.GPIO != false {
		t.Errorf("config.Mock expected (false) got (%t)", c.GPIO)
		return false
	}

	if c.MaxData != 1000 {
		t.Errorf("config.MaxData expected (1000) got (%d)", c.MaxData)
		return false
	}

	if c.Mock != false {
		t.Errorf("config.Mock expected (false) got (%t)", c.Mock)
		return false
	}

	if c.Verbose != false {
		t.Errorf("config.Mock expected (false) got (%t)", c.Verbose)
		return false
	}

	return true
}

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
