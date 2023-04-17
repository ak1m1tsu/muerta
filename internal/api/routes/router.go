package routes

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/jmoiron/sqlx"
	"github.com/romankravchuk/muerta/internal/api/middleware/notfound"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/auth"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/user"
	"github.com/romankravchuk/muerta/internal/config"
	"github.com/rs/zerolog"
)

type Router struct {
	*fiber.App
}

func NewV1(db *sqlx.DB, cfg *config.Config, logger *zerolog.Logger) *Router {
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
		r.Mount("/auth", auth.NewRouter(db))
		r.Use(jwtware.New(jwtware.Config{
			SigningMethod: "RS256",
			SigningKey:    cfg.RSAPublicKey,
		}))
		r.Mount("/users", user.NewRouter(db, logger))
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
	r.Use(csrf.New(csrf.Config{}))
}
