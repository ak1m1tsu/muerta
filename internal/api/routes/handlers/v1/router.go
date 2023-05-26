package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/auth"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/measure"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/product"
	productcategory "github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/product-category"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/recipe"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/role"
	shelflife "github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/shelf-life"
	shelflifedetector "github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/shelf-life-detector"
	shelflifestatus "github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/shelf-life-status"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/step"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/storage"
	storagetype "github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/storage-type"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/tip"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/user"
	usersetting "github.com/romankravchuk/muerta/internal/api/routes/handlers/v1/user-setting"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/redis"
)

func New(
	cfg *config.Config,
	app fiber.Router,
	client postgres.Client,
	cache redis.Client,
	logger logger.Logger,
) {
	jware := jware.New(cfg, logger)
	app.Mount("/auth", auth.NewRouter(cfg, client, logger, cache, jware))
	app.Mount("/shelf-life-detector", shelflifedetector.NewRouter(cfg, logger, jware))
	app.Mount("/recipes", recipe.NewRouter(client, logger, jware))
	app.Mount("/users", user.NewRouter(client, logger, jware))
	app.Mount("/settings", usersetting.NewRouter(client, logger, jware))
	app.Mount("/storages", storage.NewRouter(client, logger, jware))
	app.Mount("/products", product.NewRouter(client, logger, jware))
	app.Mount("/roles", role.NewRouter(client, logger, jware))
	app.Mount("/product-categories", productcategory.NewRouter(client, logger, jware))
	app.Mount("/tips", tip.NewRouter(client, logger, jware))
	app.Mount("/measures", measure.NewRouter(client, logger, jware))
	app.Mount("/steps", step.NewRouter(client, logger, jware))
	app.Mount("/shelf-lives", shelflife.NewRouter(client, logger, jware))
	app.Mount("/shelf-life-statuses", shelflifestatus.NewRouter(client, logger, jware))
	app.Mount("/storage-types", storagetype.NewRouter(client, logger, jware))
}
