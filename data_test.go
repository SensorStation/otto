package otto

import "testing"

func TestTimeSeries(t *testing.T) {
	ts := NewTimeseries("test-int")
	if ts.Label != "test-int" {
		t.Errorf("Timeseries label expected (test-int) got (%s)", ts.Label)
	}

	if len(ts.Data) != 0 {
		t.Errorf("Timeseries.Data expected len (0) got (%d)", len(ts.Data))
	}

	cnt := 10
	for i := 0; i < cnt; i++ {
		ts.Add(i)
	}

	if ts.Len() != cnt {
		t.Errorf("Timeseries.Len expected (%d) got (%d)", cnt, ts.Len())
	}

	for i := 0; i < cnt; i++ {
		d := ts.Data[i]
		if d.Int() != i {
			t.Errorf("Timeseries data expected (%d) got (%+v)", cnt, d.Int())
		}
	}
}


func TestTimeData(t *testing.T) {
	dfloat := NewData(23.3)
	val := dfloat.Float()
	if val != 23.3 {
		t.Errorf("Expected (23.3) got (%f)", val)
	}
}
