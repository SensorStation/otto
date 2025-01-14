package utils

import (
	"fmt"
	"log/slog"
	"os"
)

var (
	logfile string = "otto.log"
)

func InitLogger(lstr string, lf string) {
	if lf == "" {
		lf = logfile
	}
	level := SetLogLevel(lstr)
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("error opening log ", "err", err)
	}
	l := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(l)
}

func SetLogLevel(loglevel string) slog.Level {
	var level slog.Level

	switch loglevel {
	case "debug":
		level = slog.LevelDebug

	case "info":
		level = slog.LevelInfo

	case "warn":
		level = slog.LevelWarn

	case "error":
		level = slog.LevelError

	default:
		fmt.Printf("unknown loglevel %s sticking with warn", loglevel)
	}
	slog.SetLogLoggerLevel(level)
	return level
}
