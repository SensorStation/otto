package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sensorstation/otto/message"
)

type Websock struct {
	msgQ chan *message.Msg
	webQ map[chan *message.Msg]chan *message.Msg
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func (w *Websock) AddWebQ() chan *message.Msg {
	c := make(chan *message.Msg)
	w.webQ[c] = c
	return c
}

func (ws Websock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("[I] Connected with Websocket")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Websocket Upgrader err", "error", err)
		return
	}
	defer conn.Close()

	if ws.webQ == nil {
		ws.webQ = make(map[chan *message.Msg]chan *message.Msg)
	}

	go func() {
		for {

			println("reading a message")
			// var message StationEvent
			mt, message, err := conn.ReadMessage()
			if err != nil {
				println("read error")
				slog.Error("websocket read:", "error", err)
				break
			}
			println("read a message")
			fmt.Printf("%v - %v - %s\n", mt, message, err)
		}
	}()

	wq := ws.AddWebQ()
	for {
		msg := <-wq
		err = conn.WriteJSON(msg)
		if err != nil {
			slog.Error("Failed to write web socket", "error", err)
			return
		}
	}
}
