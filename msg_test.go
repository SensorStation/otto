package main

import (
	"fmt"
	"testing"
	"time"
)

func TestDataString(t *testing.T) {
	now := time.Now()
	d := Msg{
		Source:   "be:ef:ca:fe:01",
		Category: "data",
		Device:   "tempf",
		Value:    98.1,
		Time:     now,
	}
	formatted := fmt.Sprintf("Time: %s, Source: %s, Category: %s, Device: %s = %f",
		now.Format(time.RFC3339), d.Source, d.Category, d.Device, d.Value)

	str := d.String()
	if str != formatted {
		t.Errorf("Data Formatted expected (%s) got (%s)", formatted, str)
	}
}
