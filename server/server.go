package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
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
	return s
}

// Register to handle HTTP requests for particular paths in the
// URL or MQTT channel.
func (s *Server) Register(p string, h http.Handler) {
	slog.Info("HTTP REST API Registered: ", "path", p)
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

	slog.Info("Starting hub Web and REST server on ", "addr", s.Addr)
	http.ListenAndServe(s.Addr, s.ServeMux)
	return
}

func (s *Server) Appdir(path, file string) {
	slog.Info("appdir", "path", path)
	s.Register(path, http.FileServer(http.Dir(file)))
}

func (s *Server) EmbedTempl(path string, data any, content embed.FS) {
	slog.Info("embedTempl", "path", path)
	s.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if path == "/emb" || filepath.Ext(path) == ".html" {
			tmpl, err := template.ParseFS(content, "app/*.html")
			if err != nil {
				slog.Error("Failed to parse web template: ", "error", err.Error())
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
		slog.Error("Server.ServeHTTP failed to encode", "error", err)
	}
}
