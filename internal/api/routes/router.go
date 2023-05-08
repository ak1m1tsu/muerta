package routes

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/romankravchuk/muerta/internal/api/middleware/notfound"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/product"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/recipe"
	shelflifedetector "github.com/romankravchuk/muerta/internal/api/routes/handlers/shelf-life-detector"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/storage"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/user"
	usersetting "github.com/romankravchuk/muerta/internal/api/routes/handlers/user-setting"
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
			AppName:               "Muerta API v1.0",
			DisableStartupMessage: true,
			JSONEncoder:           sonic.Marshal,
			JSONDecoder:           sonic.Unmarshal,
		}),
	}
	r.mountAPIMiddlewares(logger)
	r.Route("/api/v1", func(r fiber.Router) {
		r.Mount("/shelf-life-detector", shelflifedetector.NewRouter(logger))
		r.Mount("/recipes", recipe.NewRouter(client, logger))
		r.Mount("/users", user.NewRouter(client, logger))
		r.Mount("/settings", usersetting.NewRouter(client, logger))
		r.Mount("/storages", storage.NewRouter(client, logger))
		r.Mount("/products", product.NewRouter(client, logger))
		// r.Mount("/auth", auth.NewRouter(cfg, db, logger))
		// r.Use(jwtware.New(jwtware.Config{
		// 	SigningMethod: "RS256",
		// 	SigningKey:    cfg.RSAPublicKey,
		// }))
	})
	r.Use(notfound.New())
	return r
}

func (r *Router) mountAPIMiddlewares(logger *log.Logger) {
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
	// r.Use(csrf.New(csrf.Config{}))
}
