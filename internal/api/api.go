package api

import (
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/routes"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
)

type API struct {
	router     *routes.Router
	listenAddr string
}

func New(client repositories.PostgresClient, cfg *config.Config, logger *log.Logger) *API {
	return &API{
		router:     routes.NewV1(client, cfg, logger),
		listenAddr: fmt.Sprintf("0.0.0.0:%s", cfg.API.Port),
	}
}

func (api *API) Run() error {
	return api.router.Listen(api.listenAddr)
}
