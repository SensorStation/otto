package logger

import (
	"log/slog"
	"os"
)

var (
	logfile string = "otto.log"
)

func init() {
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("error opening log ", "err", err)
	}
	l := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(l)
}
