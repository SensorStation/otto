package iote

import "log"

var (
	dispatcher *Dispatcher
)

func init() {
	dispatcher = NewDispatcher()
}

// Dispatcher accepts
type Dispatcher struct {
	InQ    chan *Msg
	StoreQ chan *Msg
	WebQ   map[chan *Msg]chan *Msg
}

func GetDispatcher() *Dispatcher {
	return dispatcher
}

func NewDispatcher() (d *Dispatcher) {
	d = &Dispatcher{}
	d.InQ = make(chan *Msg)
	d.WebQ = make(map[chan *Msg]chan *Msg)
	d.StoreQ = make(chan *Msg)

	go func() {

		for true {
			select {
			case msg := <-d.InQ:
				log.Printf("[I] %s", msg.String())

				switch msg.Category {
				case "d":

					for c, _ := range d.WebQ {
						c <- msg
					}
					d.StoreQ <- msg

				case "c":

				case "m":

				default:
					log.Println("Uknonwn message type: ", msg.Device)
				}
			}
		}

	}()

	return d
}

func (d *Dispatcher) AddWebQ() chan *Msg {
	c := make(chan *Msg)
	d.WebQ[c] = c
	return c
}

func (d *Dispatcher) FreeWebQ(c chan *Msg) {
	delete(d.WebQ, c)
	close(c)
}

func (d *Dispatcher) addStoreQ() chan *Msg {
	c := make(chan *Msg)
	d.StoreQ = c
	return d.StoreQ
}
