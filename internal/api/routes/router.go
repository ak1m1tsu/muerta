package routes

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/romankravchuk/muerta/internal/api/middleware/notfound"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/user"
	"github.com/rs/zerolog"
)

type Router struct {
	*fiber.App
}

func NewV1(logger *zerolog.Logger) *Router {
	r := &Router{
		App: fiber.New(fiber.Config{
			AppName:               "Muerta API v1.0",
			DisableStartupMessage: true,
			JSONEncoder:           sonic.Marshal,
			JSONDecoder:           sonic.Unmarshal,
		}),
	}
	r.mountAPIMiddlewares(logger)
	r.Route("/api/v1", func(r fiber.Router) {
		r.Mount("/users", user.NewRouter())
	})
	r.Use(notfound.New())
	return r
}

func (r *Router) mountAPIMiddlewares(logger *zerolog.Logger) {
	r.Use(requestid.New())
	r.Use(fiberzerolog.New(fiberzerolog.Config{
		Fields: []string{
			fiberzerolog.FieldRequestID,
			fiberzerolog.FieldStatus,
			fiberzerolog.FieldMethod,
			fiberzerolog.FieldPath,
			fiberzerolog.FieldLatency,
			fiberzerolog.FieldIP,
			fiberzerolog.FieldUserAgent,
			fiberzerolog.FieldError,
		},
		Logger: logger,
	}))
}
