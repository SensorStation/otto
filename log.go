package otto

import (
	"log/slog"
	"os"
)

var (
	l *Logger
)

func GetLogger() *Logger {
	return l
}

type Logger struct {
	*slog.Logger
}

func init() {
	f, err := os.OpenFile("otto.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("error opening log ", "err", err)
	}
	l = &Logger{slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))}
}
