package main

import (
	"fmt"
	"testing"
	"time"
)

func TestDataString(t *testing.T) {
	now := time.Now()
	d := Data{
		Source: "be:ef:ca:fe:01",
		Type:   "tempf",
		Value:  "98.1",
		Time:   now,
	}
	formatted := fmt.Sprintf("Time: %s, Source: %s, Type: %s = %s",
		now.Format(time.RFC3339), d.Source, d.Type, d.Value)

	str := d.String()
	if str != formatted {
		t.Errorf("Data Formatted expected (%s) got (%s)", formatted, str)
	}
}
