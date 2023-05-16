package routes

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	_ "github.com/romankravchuk/muerta/docs"
	v1 "github.com/romankravchuk/muerta/internal/api/routes/handlers/v1"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/notfound"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
)

type Router struct {
	*fiber.App
}

func NewV1(client repositories.PostgresClient, cfg *config.Config, logger *log.Logger) *Router {
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
	v1.New(cfg, routesV1, client, logger)
	r.Use(notfound.New())
	return r
}

func (r *Router) mountAPIMiddlewares(cfg *config.Config, logger *log.Logger) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Language, Content-Length, Authorization",
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
