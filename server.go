package otto

import (
	"net/http"
)

// Server serves up HTTP on Addr (default 0.0.0.0:8011)
// It takes care of REST API, serving the web app if Appdir
// does not equal nil and initial Websocket upgrade
type Server struct {
	Appdir string
	*http.Server
	*http.ServeMux
}

var (
	wserv Websock
)

func NewServer() *Server {
	s := &Server{
		Server: &http.Server{
			Addr: ":8011",
		},
	}
	s.ServeMux = http.NewServeMux()
	return s
}

// Register to handle HTTP requests for particular paths in the
// URL or MQTT channel.
func (s *Server) Register(p string, h http.Handler) {
	l.Info("HTTP REST API Registered: ", "path", p)
	s.Handle(p, h)
}

// Start the HTTP server after registering REST API callbacks
// and initializing the Web application directory
func (s *Server) Start() {
	l.Info("Starting hub Web and REST server on ", "addr", s.Addr)

	if s.Appdir != "" {
		l.Info("Server: webapp dir", "dir", s.Appdir)
		fs := http.FileServer(http.Dir(s.Appdir))
		s.Register("/", fs)
	}
	s.Register("/ws", wserv)
	s.Register("/ping", Ping{})

	l.Info("Starting HTTP server ", "addr", s.Addr)

	http.ListenAndServe(s.Addr, s.ServeMux)
	return
}
