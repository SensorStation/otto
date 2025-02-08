package devices

import "testing"

func TestDev(t *testing.T) {
	d := &BaseDevice{}
	exp := []struct {
		t string
		e bool
	}{
		{"ss/d/test", false},
	}
	for _, e := range exp {
		d.AddPub(e.t)
	}
}
