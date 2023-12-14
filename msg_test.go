package iote

import (
	"fmt"
	"testing"
	"time"
)

func TestDataString(t *testing.T) {
	now := time.Now()
	// d := MsgFloat64{
	// 	Station:  "be:ef:ca:fe:01",
	// 	Category: "data",
	// 	Device:   "tempf",
	// 	Time:     now,
	// 	Value:    98.7,
	// }

	d := MsgFloat64{}
	d.Station = "be:ef:ca:fe:01"
	d.Category = "data"
	d.Device = "tempf"
	d.Time = now
	d.Value = 98.7

	formatted := fmt.Sprintf("Time: %s, Category: %s, Station: %s, Device: %s = %f",
		now.Format(time.RFC3339), d.Category, d.Station, d.Device, d.Value)

	str := d.String()
	if str != formatted {
		t.Errorf("Data Formatted expected (%s) got (%s)", formatted, str)
	}
}
