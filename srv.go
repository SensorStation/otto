package iote

import (
	"log"
	"sync"

	"encoding/json"
	"net/http"
)

type Server struct {
	Addr   string
	Appdir string

	*http.Server
}

var (
	wserv Websock
)

// Register to handle HTTP requests for particular paths in the
// URL or MQTT channel.
func (s *Server) Register(p string, h http.Handler) {
	http.Handle(p, h)
}

// ServeHTTP provides a REST interface to the config structure
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(Stations)
	}
}

func (s *Server) Start(addr string, wg sync.WaitGroup) {
	log.Println("Starting hub Web and REST server on ", s.Addr)

	// The web app
	fs := http.FileServer(http.Dir("/srv/iot/iotvue/dist"))
	s.Register("/", fs)
	s.Register("/ws", wserv)
	s.Register("/ping", Ping{})
	s.Register("/api/data", s)
	s.Register("/api/stations", Stations)

	http.ListenAndServe(s.Addr, nil)
	wg.Done()
}
