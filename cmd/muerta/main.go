package main

import (
	"log"
	"os"

	"github.com/romankravchuk/muerta/internal/api"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	api := api.New(&logger, ":3000")
	log.Fatal(api.Run())
}
