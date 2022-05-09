package iote

import (
	"testing"
)

func TestTimeseries(t *testing.T) {
	ts := NewTimeseries()
	for i := 0; i < 100; i++ {
		ts.Append(float64(i))
	}

	l := len(ts.Values)
	if l != 100 {
		t.Errorf("expected timeseries length (100) got (%d)", l)
	}

	for i := 0.0; i < 100.0; i++ {
		v := ts.Get(int(i))
		if i != v.Val {
			t.Errorf("expected valued (%f) got (%f)", i, v.Val)
		}
	}
}

