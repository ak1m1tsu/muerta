package log

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/errors"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/rs/zerolog"
)

type Logger struct {
	zlog *zerolog.Logger
}

func New() *Logger {
	zlog := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &Logger{
		zlog: &zlog,
	}
}

func (l *Logger) GetLogger() *zerolog.Logger {
	return l.zlog
}

func (l *Logger) ServerError(ctx *fiber.Ctx, err error) {
	l.zlog.Error().Interface(
		fiberzerolog.FieldRequestID,
		ctx.GetRespHeader(fiber.HeaderXRequestID),
	).Err(err).Msg(errors.ErrServer)
}

func (l *Logger) ClientError(ctx *fiber.Ctx, err error) {
	l.zlog.Error().Interface(
		fiberzerolog.FieldRequestID,
		ctx.GetRespHeader(fiber.HeaderXRequestID),
	).Err(err).Msg(errors.ErrClient)
}

func (l *Logger) ValidationError(ctx *fiber.Ctx, errs validator.ErrorResponses) {
	l.zlog.Error().Interface(
		fiberzerolog.FieldRequestID,
		ctx.GetRespHeader(fiber.HeaderXRequestID),
	).Array(
		validator.KeyErrResponses,
		errs,
	).Msg(errors.ErrClient)
}
