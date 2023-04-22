package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/romankravchuk/muerta/internal/api"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	logger "github.com/romankravchuk/muerta/internal/pkg/log"
)

var (
	db  *sqlx.DB
	cfg *config.Config
)

func init() {
	var err error
	cfg, err = config.New(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("config create: %v", err)
	}
}

func init() {
	var err error
	db, err = sqlx.Connect("postgres", "user=postgres dbname=muerta password=postgrespw sslmode=disable")
	if err != nil {
		log.Fatalf("database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("database connection: %v", err)
	}
}

func main() {
	logger := logger.New()
	api := api.New(db, cfg, logger)
	log.Fatalf("api run: %v", api.Run())
}
