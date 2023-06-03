package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// Data is a general structure that holds a single data item
// such as a value that is read from a sensor
type Data struct {
	Source string      `json:Source`
	Type   string      `json:Source`
	Value  interface{} `json:Value`

	time.Time `json:Time`
}

func (d Data) String() string {

	var str string
	switch v := d.Value.(type) {
	case int:
		str = strconv.Itoa(v)
	case string:
		str = string(v)
	case []uint8:
		str = d.Value.(string)

	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

	str = fmt.Sprintf("Time: %s, Source: %s, Type: %s = %s",
		d.Time.Format(time.RFC3339), d.Source, d.Type, str)
	return str
}

func startDataQ() (dataQ chan *Data) {
	dataQ = make(chan *Data)

	go func() {

		for true {
			select {
			case data := <-dataQ:
				log.Printf("[I] %s", data.String())
			}
		}

	}()

	return dataQ
}
