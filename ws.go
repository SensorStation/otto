package otto

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Websock struct {
	msgQ chan *Station
	webQ map[chan *Station]chan *Station
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func (w *Websock) AddWebQ() chan *Station {
	c := make(chan *Station)
	w.webQ[c] = c
	return c
}

func (ws Websock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.Info("[I] Connected with Websocket")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.Error("Websocket Upgrader err", "error", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			var message StationEvent
			err := conn.ReadJSON(&message)
			if err != nil {
				l.Error("Websocket ", "message", message, "error", err)
				break
			}

			switch message.Type {
			case "relay":
				// Stations.EventQ <- &message

			default:
				l.Error("unknown event type", "message", message)
			}

		}
	}()

	wq := ws.AddWebQ()
	for {
		msg := <-wq
		err = conn.WriteJSON(msg)
		if err != nil {
			l.Error("Failed to write web socket", "error", err)
			return
		}
	}
}
