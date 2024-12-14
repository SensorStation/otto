package otto

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// Server serves up HTTP on Addr (default 0.0.0.0:8011)
// It takes care of REST API, serving the web app if Appdir
// does not equal nil and initial Websocket upgrade
type Server struct {
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

	s.Register("/ws", wserv)
	s.Register("/ping", Ping{})
	l.Info("Starting HTTP server ", "addr", s.Addr)
	http.ListenAndServe(s.Addr, s.ServeMux)
	return
}

func (s *Server) Appdir(path, file string) {
	s.Register(path, http.FileServer(http.Dir(file)))
}

func (s *Server) EmbedTempl(path string, data any, content embed.FS) {
	fmt.Println("PATH: ", path)

	s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if path == "/" || filepath.Ext(path) == ".html" {

			fmt.Println("here we are foo bar ", r.URL.String())
			tmpl, err := template.ParseFS(content, "app/*.html")
			if err != nil {
				l.Error("Failed to parse web template: ", "error", err.Error())
			}

			tmpl.Execute(w, data)

		} else {
			s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("there we go ", r.URL.String())
			})
		}

	})
}
