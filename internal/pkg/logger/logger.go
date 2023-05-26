// Package log provides a structured logging solution for the application.
// It utilizes the zerolog library to log messages to stderr.
package logger

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

const (
	errClient = "Client Error"
	errServer = "Server Error"
)

type Type int

const (
	Client Type = iota + 1
	Server
	Validation
)

type Logger interface {
	Error(*fiber.Ctx, Type, error)
	GetLogger() *zerolog.Logger
}

// Logger is a struct that contains the zerolog logger.
type logger struct {
	zlog *zerolog.Logger
}

// New creates a new Logger instance with a new zerolog logger.
func New() Logger {
	zlog := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &logger{
		zlog: &zlog,
	}
}

// Error logs an error message with the associated request ID and error.
func (l *logger) Error(ctx *fiber.Ctx, t Type, err error) {
	msg := ""
	switch t {
	case Client | Validation:
		msg = errClient
	case Server:
		msg = errServer
	}
	l.zlog.Error().Interface(
		fiberzerolog.FieldRequestID,
		ctx.GetRespHeader(fiber.HeaderXRequestID),
	).Err(err).Msg(msg)
}

// GetLogger returns the zerolog logger contained within the Logger.
func (l *logger) GetLogger() *zerolog.Logger {
	return l.zlog
}
