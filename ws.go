package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WSServer struct {
	C *websocket.Conn
	RecvQ chan Msg
}

type KeyVal struct {
	K string
	V interface{}
}

var (
	WServ WSServer
)

func (ws WSServer) GetID() string {
	return "Websocket Server"
}

func (ws WSServer) GetRecvQ() chan Msg {
	return ws.RecvQ
}

// ServeHTTP
func (ws WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Websocket connection ")
	var err error
	ws.C, err = websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:       []string{"echo"},
		InsecureSkipVerify: true, // Take care of CORS
		// OriginPatterns: ["*"],
	})

	ws.RecvQ = make(chan Msg)
	if err != nil {
		log.Println("ERROR ", err)
		return
	}
	defer ws.C.Close(websocket.StatusInternalError, "houston, we have a problem")
	hub.AddConsumer("data", ws)

	running := true
	go func() {
		for running {
			select {
			case msg := <-ws.RecvQ:
				m64 := msg.ToMsgFloat64()
				err = wsjson.Write(r.Context(), ws.C, m64)
				/*
				err = wsjson.Write(r.Context(), ws.C, map[string]string{
					"station": msg.Station,
					"sensor": msg.Sensor,
					"time": msg.Time,
					"value": val,
				})
				if err != nil {
					log.Println("Error sending websock: ", err)
				}
				*/
			}
		}
	}()

	for running {
		data := make([]byte, 8192)
		_, data, err := ws.C.Read(r.Context())
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			log.Println("ws Closed")
			return
		}
		if err != nil {
			log.Println("ERROR: reading websocket ", err)
			return
		}
		log.Printf("incoming: %s", string(data))
	}

}

func echo(ctx context.Context, c *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}
