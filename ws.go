package otto

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Websock struct {
	msgQ chan *Station
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func (ws Websock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[I] Connected with Websocket")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket Upgrader err", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			var message StationEvent
			err := conn.ReadJSON(&message)
			if err != nil {
				log.Println("read:", err)
				break
			}

			switch message.Type {
			case "relay":
				Stations.EventQ <- &message

			default:
				log.Printf("ERROR: unknown event type: %+v", message)
			}

		}
	}()

	wq := o.Dispatcher.AddWebQ()
	defer o.Dispatcher.FreeWebQ(wq)

	for {
		msg := <-wq
		err = conn.WriteJSON(msg)
		if err != nil {
			log.Println("Failed to write web socket", err)
			return
		}
	}
}
