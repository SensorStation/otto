package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/sensorstation/otto/logger"
)

// Server serves up HTTP on Addr (default 0.0.0.0:8011)
// It takes care of REST API, serving the web app if Appdir
// does not equal nil and initial Websocket upgrade
type Server struct {
	*http.Server
	*http.ServeMux

	EndPoints map[string]http.Handler
}

var (
	wserv  Websock
	server *Server
	l      *logger.Logger
)

func GetServer() *Server {
	if server == nil {
		server = NewServer()
	}
	return server
}

func NewServer() *Server {
	s := &Server{
		Server: &http.Server{
			Addr: ":8011",
		},
	}
	s.ServeMux = http.NewServeMux()
	if l == nil {
		l = logger.GetLogger()
	}
	return s
}

// Register to handle HTTP requests for particular paths in the
// URL or MQTT channel.
func (s *Server) Register(p string, h http.Handler) {
	l.Info("HTTP REST API Registered: ", "path", p)
	if s.EndPoints == nil {
		s.EndPoints = make(map[string]http.Handler)
	}
	s.EndPoints[p] = h
	s.Handle(p, h)
}

// Start the HTTP server after registering REST API callbacks
// and initializing the Web application directory
func (s *Server) Start() {
	s.Register("/ws", wserv)
	s.Register("/ping", Ping{})
	s.Register("/api", s)

	l.Info("Starting hub Web and REST server on ", "addr", s.Addr)
	http.ListenAndServe(s.Addr, s.ServeMux)
	return
}

func (s *Server) Appdir(path, file string) {
	l.Info("appdir", "path", path)
	s.Register(path, http.FileServer(http.Dir(file)))
}

func (s *Server) EmbedTempl(path string, data any, content embed.FS) {
	l.Info("embedTempl", "path", path)
	s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if path == "/emb" || filepath.Ext(path) == ".html" {
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

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ep := struct {
		Endpoints []string
	}{}
	for e, _ := range s.EndPoints {
		ep.Endpoints = append(ep.Endpoints, e)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ep)
	if err != nil {
		l.Error("Server.ServeHTTP failed to encode", "error", err)
	}
}
