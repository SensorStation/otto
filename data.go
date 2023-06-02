package main

import (
	"fmt"
	"time"
)

// Data is a general structure that holds a single data item
// such as a value that is read from a sensor
type Data struct {
	Source	string		`json:Source`
	Type	string		`json:Source`
	Value	interface{} `json:Value`

	time.Time		    `json:Time`
}

func (d Data) String() string {
	str := fmt.Sprintf("Time: %s: Source: %s, Type: %s = %s\n",
		d.Time.Format(time.RFC3339), d.Source, d.Type, d.Value.([]uint8))
	return str
}
