package otto

import "testing"

func TestStats(t *testing.T) {
	st := GetStats()

	if st.Goroutines == 0 {
		t.Error("Surely there must be more than one goroutine")
	}

	if st.CPUs == 0 {
		t.Error("Did not fund any CPUs")
	}

	if st.GoVersion == "" {
		t.Error("Did not get version")
	}
}
