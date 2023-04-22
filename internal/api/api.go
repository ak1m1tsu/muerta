package api

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/romankravchuk/muerta/internal/api/routes"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
)

type API struct {
	router     *routes.Router
	listenAddr string
}

func New(db *sqlx.DB, cfg *config.Config, logger *log.Logger) *API {
	return &API{
		router:     routes.NewV1(db, cfg, logger),
		listenAddr: fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port),
	}
}

func (api *API) Run() error {
	return api.router.Listen(api.listenAddr)
}
