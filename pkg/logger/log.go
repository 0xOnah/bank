package logger

import (
	"log/slog"
	"os"
)

func init() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}).WithAttrs([]slog.Attr{})

	log := slog.New(logHandler)
	slog.SetDefault(log)
}
