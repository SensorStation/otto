package main

import (
	"fmt"
	"log"
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

func startMsgQ() (msgQ chan *Msg) {
	msgQ = make(chan *Msg)

	go func() {

		for true {
			select {
			case msg := <-msgQ:
				log.Printf("[I] %s", msg.String())

				src := msg.Source
				switch msg.Category {
				case "data":

					// if there are websockets waiting to recieve this
					// data send it to them

					store.Store(msg)

				case "control":
					log.Println("Do something with the control from ", src)

				default:
					log.Println("Uknonwn message type: ", msg.Device)
				}

			}
		}

	}()

	return msgQ
}
