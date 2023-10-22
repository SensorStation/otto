package main

import (
	"fmt"
	"strconv"
	"time"
)

// Data is a general structure that holds a single data item
// such as a value that is read from a sensor
type Msg struct {
	Source   string      `json:source`
	Category string      `json:category`
	Device   string      `json:device`
	Value    interface{} `json:value`

	time.Time `json:Time`
}

func (d Msg) String() string {

	var str string
	switch v := d.Value.(type) {
	case int:
		str = strconv.Itoa(v)
	case string:
		str = string(v)
	case []uint8:
		for _, c := range v {
			str += string(c)
		}

	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

	str = fmt.Sprintf("Source: %s, Device: %s = %s", d.Source, d.Device, str)
	return str
}
