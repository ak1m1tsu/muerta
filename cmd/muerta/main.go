package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

// main start point of the application
//
//	@title						Muerta API
//	@version					1.0.0
//	@description				Web API to control the shelf life of products using computer vision
//	@termsOfService				http://swagger.io/terms
//
//	@BasePath					/api/v1
//
//	@securitydefinitions.apiKey	BearerAuth
//	@in							header
//	@name						Authrization
func main() {
	logger := logger.New()
	api := api.New(client, cfg, logger)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		cfg.ShutdownShelfDetectorChan <- struct{}{}
		_ = api.Shutdown()
	}()
	log.Fatalf("api run: %v", api.Run())
}
