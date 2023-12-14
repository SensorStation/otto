package iote

import (
	"log"
	"sync"

	"encoding/json"
	"net/http"
)

type Server struct {
	Addr string
	*http.Server
}

func NewServer(addr string) *Server {
	return &Server{Addr: addr}
}

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
	http.ListenAndServe(s.Addr, nil)
	wg.Done()
}
