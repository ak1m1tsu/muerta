// Package log provides a structured logging solution for the application.
// It utilizes the zerolog library to log messages to stderr.
package log

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/rs/zerolog"
)

const (
	errClient = "Client Error"
	errServer = "Server Error"
)

// Logger is a struct that contains the zerolog logger.
type Logger struct {
	zlog *zerolog.Logger
}

// New creates a new Logger instance with a new zerolog logger.
func New() *Logger {
	zlog := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &Logger{
		zlog: &zlog,
	}
}

// GetLogger returns the zerolog logger contained within the Logger.
func (l *Logger) GetLogger() *zerolog.Logger {
	return l.zlog
}

// ServerError logs a server error message with the associated request ID and error.
func (l *Logger) ServerError(ctx *fiber.Ctx, err error) {
	l.zlog.Error().Interface(
		fiberzerolog.FieldRequestID,
		ctx.GetRespHeader(fiber.HeaderXRequestID),
	).Err(err).Msg(errServer)
}

// ClientError logs a client error message with the associated request ID and error.
func (l *Logger) ClientError(ctx *fiber.Ctx, err error) {
	l.zlog.Error().Interface(
		fiberzerolog.FieldRequestID,
		ctx.GetRespHeader(fiber.HeaderXRequestID),
	).Err(err).Msg(errClient)
}

// ValidationError logs a validation error message with the associated request ID and errors.
func (l *Logger) ValidationError(ctx *fiber.Ctx, errs validator.ValidationErrors) {
	l.zlog.Error().Interface(
		fiberzerolog.FieldRequestID,
		ctx.GetRespHeader(fiber.HeaderXRequestID),
	).Array(
		validator.KeyErrResponses,
		errs,
	).Msg(errClient)
}
