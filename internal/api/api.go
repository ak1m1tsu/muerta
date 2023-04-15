package api

import (
	"github.com/romankravchuk/muerta/internal/api/routes"
	"github.com/rs/zerolog"
)

type API struct {
	router *routes.Router
	port   string
}

func New(logger *zerolog.Logger, port string) *API {
	return &API{
		router: routes.NewV1(logger),
		port:   port,
	}
}

func (api *API) Run() error {
	return api.router.Listen(api.port)
}
