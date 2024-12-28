package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sensorstation/otto/station"
)

type Websock struct {
	msgQ chan *station.Station
	webQ map[chan *station.Station]chan *station.Station
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func (w *Websock) AddWebQ() chan *station.Station {
	c := make(chan *station.Station)
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

	if ws.webQ == nil {
		ws.webQ = make(map[chan *station.Station]chan *station.Station)
	}

	go func() {
		for {

			println("reading a message")
			// var message StationEvent
			mt, message, err := conn.ReadMessage()
			if err != nil {
				println("read error")
				l.Error("websocket read:", "error", err)
				break
			}
			println("read a message")
			fmt.Printf("%v - %v - %s\n", mt, message, err)
			// if err != nil {
			// 	l.Error("Websocket ", "message", message, "error", err)
			// 	break
			// }

			// switch message.Type {
			// case "relay":
			// 	// Stations.EventQ <- &message

			// default:
			// 	l.Error("unknown event type", "message", message)
			// }

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
