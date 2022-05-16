package main

import "testing"

func TestHub(t *testing.T) {
	st := NewHub(&config)
	if st.Addr != "0.0.0.0:8011" {
		t.Errorf("hub Addr incorrect Expected (0.0.0.0:8011) got (%s)", st.Addr)
	}
}
