package ctxlog

import (
	"log/slog"
	"os"
)

var LogLevel = slog.LevelInfo
var Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: LogLevel,
}))

func ParseLevel(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func SetupLogger(levelStr string) {
	LogLevel = ParseLevel(levelStr)
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: LogLevel,
	}))
}
