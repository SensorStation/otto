package main

import (
	"embed"

	"github.com/sensorstation/otto"
)

//go:embed app
var content embed.FS

func initApp() {
	s := otto.GetServer()

	// The following line is commented out because
	var data any
	s.EmbedTempl("/emb", data, content)
	s.Appdir("/", "app")
	go s.Start()
}
