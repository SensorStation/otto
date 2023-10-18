package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type websock struct {
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ws websock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// register this connection for data recieved by websocket

	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			return
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			return
		}
		if _, err := io.Copy(w, r); err != nil {
			return
		}
		if err := w.Close(); err != nil {
			return
		}
	}
}
