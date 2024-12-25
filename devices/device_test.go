package devices

import "testing"

func TestDev(t *testing.T) {
	d := &Dev{}
	name := "test-dev"
	d.SetName(name)
	if d.Name() != name {
		t.Errorf("Bad name expected (%s) got (%s)", name, d.Name())
	}

	exp := []struct {
		t string
		e bool
	}{
		{"ss/c/test", false},
		{"ss/d/test", false},
	}
	for _, e := range exp {
		d.AddPub(e.t)
	}

	for _, p := range d.Pubs() {
		for i, _ := range exp {
			if exp[i].t == p {
				exp[i].e = true
				break
			}
		}
	}

	for _, e := range exp {
		if !e.e {
			t.Errorf("expected to find topic (%s) but did not", e.t)
		}
	}

	if d.Period() != 0 {
		t.Errorf("period expected (0) got (%d)", d.Period())
	}
}
