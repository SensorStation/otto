package otto

import (
	"log"

	"encoding/json"
	"net/http"
)

// Server serves up HTTP on Addr (default 0.0.0.0:8011)
// It takes care of REST API, serving the web app if Appdir
// does not equal nil and initial Websocket upgrade
type Server struct {
	Addr   string
	Appdir string
}

var (
	wserv Websock
)

// Register to handle HTTP requests for particular paths in the
// URL or MQTT channel.
func (s *Server) Register(p string, h http.Handler) {
	log.Println("HTTP REST API Registered: ", p)
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

// Start the HTTP server after registering REST API callbacks
// and initializing the Web application directory
func (s *Server) Start() {
	log.Println("Starting hub Web and REST server on ", s.Addr)

	if s.Appdir != "" {
		log.Println("Server: webapp dir", s.Appdir)
		fs := http.FileServer(http.Dir(s.Appdir))
		s.Register("/", fs)
	}
	s.Register("/ws", wserv)
	s.Register("/ping", Ping{})
	s.Register("/api/data", s)
	s.Register("/api/stations", Stations)

	go http.ListenAndServe(s.Addr, nil)
	return
}
