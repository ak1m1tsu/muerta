package logger

import (
	"log/slog"
	"os"
)

const (
	dev   = "dev"
	prod  = "prod"
	local = "local"
)

func New(env string) *slog.Logger {
	var s *slog.Logger
	switch env {
	case dev:
		s = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prod:
		s = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case local:
		s = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		s = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}
	return s
}
