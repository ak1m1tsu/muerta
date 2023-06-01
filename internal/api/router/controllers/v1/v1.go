package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/auth"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/measure"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/product"
	productcategory "github.com/romankravchuk/muerta/internal/api/router/controllers/v1/product-category"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/recipe"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/role"
	shelflife "github.com/romankravchuk/muerta/internal/api/router/controllers/v1/shelf-life"
	shelflifedetector "github.com/romankravchuk/muerta/internal/api/router/controllers/v1/shelf-life-detector"
	shelflifestatus "github.com/romankravchuk/muerta/internal/api/router/controllers/v1/shelf-life-status"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/step"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/tip"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/user"
	usersetting "github.com/romankravchuk/muerta/internal/api/router/controllers/v1/user-setting"
	"github.com/romankravchuk/muerta/internal/api/router/controllers/v1/vault"
	storagetype "github.com/romankravchuk/muerta/internal/api/router/controllers/v1/vaulttype"
	jware "github.com/romankravchuk/muerta/internal/api/router/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/redis"
)

func New(
	cfg *config.Config,
	app fiber.Router,
	db postgres.Client,
	cache redis.Client,
	log logger.Logger,
) {
	jware := jware.New(cfg, log)
	app.Mount("/auth", auth.NewRouter(cfg, db, log, cache, jware))
	app.Mount("/shelf-life-detector", shelflifedetector.NewRouter(cfg, log, jware))
	app.Mount("/recipes", recipe.NewRouter(db, log, jware))
	app.Mount("/users", user.NewRouter(db, log, jware))
	app.Mount("/settings", usersetting.NewRouter(db, log, jware))
	app.Mount("/storages", vault.NewRouter(db, log, jware))
	app.Mount("/products", product.NewRouter(db, log, jware))
	app.Mount("/roles", role.NewRouter(db, log, jware))
	app.Mount("/product-categories", productcategory.NewRouter(db, log, jware))
	app.Mount("/tips", tip.NewRouter(db, log, jware))
	app.Mount("/measures", measure.NewRouter(db, log, jware))
	app.Mount("/steps", step.NewRouter(db, log, jware))
	app.Mount("/shelf-lives", shelflife.NewRouter(db, log, jware))
	app.Mount("/shelf-life-statuses", shelflifestatus.NewRouter(db, log, jware))
	app.Mount("/storage-types", storagetype.NewRouter(db, log, jware))
}
