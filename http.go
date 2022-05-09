package iote

import (
	"encoding/json"
	"net/http"
)

var (
	Server *HTTPServer = &HTTPServer{}
)

type HTTPServer struct {
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(stations)
	}
}

// Register to handle HTTP requests for particular paths in the
// URL or MQTT channel.
func (s *HTTPServer) Register(p string, h http.Handler) {
	http.Handle(p, h)
}

func (s *HTTPServer) Listen() error {
	return http.ListenAndServe(config.Addr, nil)
}
