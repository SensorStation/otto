package relay

import "testing"

func TestRelay(t *testing.T) {
	r := New("relay", 5)
	if r.Name() != "relay" {
		t.Errorf("relay expected Name (%s) got (%s)", "relay", r.Name())
	}
}
