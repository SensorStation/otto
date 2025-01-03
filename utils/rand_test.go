package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRando(t *testing.T) {
	r := NewRando()
	v := r.Float64()
	if v <= 0.0 {
		t.Errorf("Expected a float value got (%f2.1)", v)
	}

	str := r.String()
	if str == "" {
		t.Errorf("expected a float value got (%s)", r.String())
	}
}

func TestRandHTTP(t *testing.T) {
	ts := httptest.NewServer(NewRando())
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
	fmt.Printf("%+v\n", cbuf)
}
