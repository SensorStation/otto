package main

import (
	"fmt"
	"time"
)

// Data is a general structure that holds a single data item
// such as a value that is read from a sensor
type Msg struct {
	Source   string  `json:"source"`
	Category string  `json:"category"`
	Device   string  `json:"device"`
	Value    float64 `json:"value"`

	Time time.Time `json:"time"`
}

func (d Msg) String() string {

	var str string
	str = fmt.Sprintf("Time: %s, Source: %s, Category: %s, Device: %s = %f",
		d.Time.Format(time.RFC3339),
		d.Source, d.Category, d.Device, d.Value)
	return str
}
