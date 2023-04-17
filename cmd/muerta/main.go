package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/romankravchuk/muerta/internal/api"
	"github.com/romankravchuk/muerta/internal/config"
	"github.com/rs/zerolog"
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
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	api := api.New(db, cfg, &logger)
	log.Fatalf("api run: %v", api.Run())
}
