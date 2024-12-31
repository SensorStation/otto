package devices

import "testing"

func TestDev(t *testing.T) {
	d := &Device{}
	exp := []struct {
		t string
		e bool
	}{
		{"ss/d/test", false},
	}
	for _, e := range exp {
		d.AddPub(e.t)
	}

	for i, _ := range exp {
		if exp[i].t == d.Pub {
			exp[i].e = true
			break
		}
	}

	for _, e := range exp {
		if !e.e {
			t.Errorf("expected to find topic (%s) but did not", e.t)
		}
	}

	if d.Period != 0 {
		t.Errorf("period expected (0) got (%d)", d.Period)
	}
}
