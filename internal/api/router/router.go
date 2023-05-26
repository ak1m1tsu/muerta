package router

import (
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	_ "github.com/romankravchuk/muerta/internal/api/docs"
	v1 "github.com/romankravchuk/muerta/internal/api/router/controllers/v1"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/notfound"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/redis"
)

type Router struct {
	*fiber.App
}

func NewV1(
	cfg *config.Config,
	client postgres.Client,
	cache redis.Client,
	logger logger.Logger,
) *Router {
	r := &Router{
		App: fiber.New(fiber.Config{
			AppName:     "Muerta API v1.0",
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		}),
	}
	r.mountAPIMiddlewares(cfg, logger)
	r.Get("/docs/*", swagger.HandlerDefault)
	api := r.Group("/api")
	routesV1 := api.Group("/v1")
	v1.New(cfg, routesV1, client, cache, logger)
	r.Use(notfound.New())
	return r
}

func (r *Router) mountAPIMiddlewares(cfg *config.Config, logger logger.Logger) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Language, Content-Length, Authorization",
	}))
	r.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/": "/docs",
		},
		StatusCode: http.StatusMovedPermanently,
	}))
	r.Use(requestid.New())
	r.Use(recover.New())
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
		Logger: logger.GetLogger(),
	}))
}
