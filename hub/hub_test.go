package hub

import "testing"

func TestHub(t *testing.T) {
	hub := NewHub(&config)
	t.Logf("HUB: %+v", hub)
}
