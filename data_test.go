package otto

import "testing"

func TestTimeData(t *testing.T) {
	dfloat := NewData(23.3)
	val := dfloat.Float()
	if val != 23.3 {
		t.Errorf("Expected (23.3) got (%f)", val)
	}
}
