package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type websock struct {
	msgQ chan *Msg
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func (ws websock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("First err", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}()

	wq := disp.addWebQ()
	defer disp.freeWebQ(wq)

	for {
		msg := <-wq
		err = conn.WriteJSON(msg)
		if err != nil {
			log.Println("Failed to write web socket", err)
			return
		}
	}

	log.Println("WS Connection going to close")
}
