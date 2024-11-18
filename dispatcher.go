package otto

// import (
// 	"log"
// )

// // Dispatcher accepts
// type Dispatcher struct {
// 	InQ    chan *Msg
// 	StoreQ chan *Msg
// 	WebQ   map[chan *Station]chan *Station
// }

// func NewDispatcher() (d *Dispatcher) {
// 	d = &Dispatcher{}
// 	d.InQ = make(chan *Msg)
// 	d.StoreQ = make(chan *Msg)
// 	d.WebQ = make(map[chan *Station]chan *Station)

// 	go func() {

// 		for true {

// 			select {
// 			case msg := <-d.InQ:
// 				switch msg.Type {
// 				case "station":

// 					st := Stations.Update(msg)
// 					for c, _ := range d.WebQ {
// 						c <- st
// 					}

// 				default:
// 					log.Println("Uknonwn message type: ", msg.Type)
// 				}
// 			}
// 		}
// 	}()

// 	return d
// }

// func (d *Dispatcher) FreeWebQ(c chan *Station) {
// 	delete(d.WebQ, c)
// 	close(c)
// }

// func (d *Dispatcher) addStoreQ() chan *Msg {
// 	c := make(chan *Msg)
// 	d.StoreQ = c
// 	return d.StoreQ
// }
