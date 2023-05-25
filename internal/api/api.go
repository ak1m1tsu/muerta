package api

import (
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/storage/redis"
)

type API struct {
	router     *routes.Router
	listenAddr string
}

func New(
	cfg *config.Config,
	client repositories.PostgresClient,
	cache redis.Client,
	logger *log.Logger,
) *API {
	return &API{
		router:     routes.NewV1(cfg, client, cache, logger),
		listenAddr: fmt.Sprintf("0.0.0.0:%s", cfg.API.Port),
	}
}

func (api *API) Run() error {
	return api.router.Listen(api.listenAddr)
}

func (api *API) Shutdown() error {
	return api.router.Shutdown()
}
