package otto

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	ts := httptest.NewServer(Ping{})
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		l.Fatal(err)
	}
	pong, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		l.Fatal(err)
	}

	if string(pong) != "Pong\n" {
		l.Fatal(err)
	}
}
