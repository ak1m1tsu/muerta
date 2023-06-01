package api

import (
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/redis"
)

type API struct {
	router     *router.Router
	listenAddr string
}

func New(
	cfg *config.Config,
	client postgres.Client,
	cache redis.Client,
	logger logger.Logger,
) *API {
	return &API{
		router:     router.NewV1(cfg, client, cache, logger),
		listenAddr: fmt.Sprintf("0.0.0.0:%s", cfg.API.Port),
	}
}

func (api *API) Run() error {
	return api.router.Listen(api.listenAddr)
}

func (api *API) Shutdown() error {
	return api.router.Shutdown()
}
