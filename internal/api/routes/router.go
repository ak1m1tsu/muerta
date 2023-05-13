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
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/auth"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/measure"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/product"
	productcategory "github.com/romankravchuk/muerta/internal/api/routes/handlers/product-category"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/recipe"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/role"
	shelflife "github.com/romankravchuk/muerta/internal/api/routes/handlers/shelf-life"
	shelflifedetector "github.com/romankravchuk/muerta/internal/api/routes/handlers/shelf-life-detector"
	shelflifestatus "github.com/romankravchuk/muerta/internal/api/routes/handlers/shelf-life-status"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/step"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/storage"
	storagetype "github.com/romankravchuk/muerta/internal/api/routes/handlers/storage-type"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/tip"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/user"
	usersetting "github.com/romankravchuk/muerta/internal/api/routes/handlers/user-setting"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
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
	jware := jware.New(cfg, logger)
	r.mountAPIMiddlewares(logger)
	r.Route("/api/v1", func(r fiber.Router) {
		r.Get("/swagger/*", swagger.HandlerDefault)
		r.Mount("/shelf-life-detector", shelflifedetector.NewRouter(logger, jware))
		r.Mount("/recipes", recipe.NewRouter(client, logger, jware))
		r.Mount("/users", user.NewRouter(client, logger, jware))
		r.Mount("/settings", usersetting.NewRouter(client, logger, jware))
		r.Mount("/storages", storage.NewRouter(client, logger, jware))
		r.Mount("/products", product.NewRouter(client, logger, jware))
		r.Mount("/roles", role.NewRouter(client, logger, jware))
		r.Mount("/product-categories", productcategory.NewRouter(client, logger, jware))
		r.Mount("/tips", tip.NewRouter(client, logger, jware))
		r.Mount("/measures", measure.NewRouter(client, logger, jware))
		r.Mount("/steps", step.NewRouter(client, logger, jware))
		r.Mount("/shelf-lives", shelflife.NewRouter(client, logger, jware))
		r.Mount("/shelf-life-statuses", shelflifestatus.NewRouter(client, logger, jware))
		r.Mount("/storage-types", storagetype.NewRouter(client, logger, jware))
		r.Mount("/auth", auth.NewRouter(cfg, client, logger, jware))
	})
	r.Use(notfound.New())
	return r
}

func (r *Router) mountAPIMiddlewares(logger *log.Logger) {
	r.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowCredentials: true,
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
