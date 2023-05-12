package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/romankravchuk/muerta/internal/api"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	logger "github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
)

var (
	client *pgxpool.Pool
	cfg    *config.Config
)

func init() {
	var err error
	cfg, err = config.New()
	if err != nil {
		log.Fatalf("config create: %v", err)
	}
}

func init() {
	var err error
	client, err = repositories.NewPostgresClient(context.Background(), 5, cfg)
	if err != nil {
		log.Fatalf("database connection: %v", err)
	}
}

//	@title						Muerta API
//	@version					1.0
//	@description				API for Muerta
//	@BasePath					/api/v1
//	@securityDefinitions.jwt	BearerAuth
func main() {
	logger := logger.New()
	api := api.New(client, cfg, logger)
	log.Fatalf("api run: %v", api.Run())
}
