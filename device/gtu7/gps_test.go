package gtu7

import (
	"testing"
)

func TestGPS(t *testing.T) {
	// need to add more tests
	gotone := false
	g := NewGTU7("fakedev")
	g.OpenStrings(input)
	gpsQ := g.StartReading()
	for gps := range gpsQ {
		if gps.complete {
			gotone = true
		}
	}
	if !gotone {
		t.Error("Failed to recieve a complete GPS package")
	}
}

var input string = `
$GPGGA,160446.00,3340.34121,N,11800.11332,W,2,08,1.20,11.8,M,-33.1,M,,0000*58
$GPGSA,A,3,09,16,46,03,07,31,26,04,,,,,3.08,1.20,2.84*0E
$GPGSV,4,1,13,01,02,193,,03,58,181,33,04,64,360,31,06,12,295,*7A
$GPGSV,4,2,13,07,32,254,25,08,00,154,,09,44,317,33,16,52,085,26*72
$GPGSV,4,3,13,26,31,051,15,27,05,124,16,31,15,053,10,46,49,200,33*76
$GPGSV,4,4,13,48,50,193,*49
$GPGLL,3340.34121,N,11800.11332,W,160446.00,A,D*74
$GPRMC,160447.00,A,3340.34118,N,11800.11331,W,0.063,,020125,,,D*64
$GPVTG,,T,,M,0.063,N,0.117,K,D*24
`
