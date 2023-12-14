package iote

import (
	"log"
)

var (
	disp Dispatcher
)

// Dispatcher accepts
type Dispatcher struct {
	InQ    chan *Msg
	StoreQ *chan *Msg
	webQ   map[chan *Msg]chan *Msg
}

func (d *Dispatcher) AddWebQ() chan *Msg {
	c := make(chan *Msg)
	d.webQ[c] = c
	return c
}

func (d *Dispatcher) FreeWebQ(c chan *Msg) {
	delete(d.webQ, c)
	close(c)
}

func (d *Dispatcher) addStoreQ() *chan *Msg {
	c := make(chan *Msg)
	d.StoreQ = &c
	return d.StoreQ
}

func NewDispatcher() (d *Dispatcher) {
	d = &Dispatcher{}
	d.InQ = make(chan *Msg)
	d.webQ = make(map[chan *Msg]chan *Msg)

	go func() {

		for true {
			select {
			case msg := <-d.InQ:
				log.Printf("[I] %s", msg.String())

				switch msg.Category {
				case "data":

					// if there are websockets waiting to recieve this
					// data send it to them
					if d.StoreQ != nil {
						*d.StoreQ <- msg
					}

					for c, _ := range d.webQ {
						c <- msg
					}

				case "control":

				default:
					log.Println("Uknonwn message type: ", msg.Device)
				}

			}
		}

	}()

	return d
}
