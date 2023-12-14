package main

import "log"

// Dispatcher accepts
type dispatcher struct {
	InQ    chan *Msg
	StoreQ *chan *Msg
	webQ   map[chan *Msg]chan *Msg
}

func (d *dispatcher) addWebQ() chan *Msg {
	c := make(chan *Msg)
	d.webQ[c] = c
	return c
}

func (d *dispatcher) freeWebQ(c chan *Msg) {
	delete(d.webQ, c)
	close(c)
}

func (d *dispatcher) addStoreQ() *chan *Msg {
	c := make(chan *Msg)
	d.StoreQ = &c
	return d.StoreQ
}

func newDispatcher() (d *dispatcher) {
	d = &dispatcher{}
	d.InQ = make(chan *Msg)
	d.webQ = make(map[chan *Msg]chan *Msg)

	go func() {

		for true {
			select {
			case msg := <-d.InQ:
				log.Printf("[I] %s", msg.String())

				src := msg.Source
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
					log.Println("Do something with the control from ", src)

				default:
					log.Println("Uknonwn message type: ", msg.Device)
				}

			}
		}

	}()

	return d
}
