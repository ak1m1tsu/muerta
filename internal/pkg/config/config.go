// Package config provides a configuration struct and a constructor function
// for initializing the configuration object with values from environment
// variables and certificate files.
//
// Example usage:
//
//	cfg, err := config.New()
//	if err != nil {
//	    log.Fatalf("Failed to load configuration: %v", err)
//	}
//	fmt.Printf("API name: %s, Port: %s", cfg.API.Name, cfg.API.Port)
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// The Config struct contains fields for the API and database connection
// settings, as well as access and refresh token configurations.
type Config struct {
	API struct {
		// Name of the API service
		Name string
		// Port number to listen on
		Port string
	}
	Database struct {
		// Host name of the database server
		Host string
		// Port number of the database server
		Port string
		// Username for the database authentication
		User string
		// Password for the database authentication
		Password string
		// Name of the database to connect to
		Name string
	}
	Cache struct {
		// Host name of the redis server
		Host string
		// Port number of the redis server
		Port string
		// Username for the redis authentication
		User string
		// Password for the redis authentication
		Password string
	}
	// Private key for signing access tokens
	AccessTokenPrivateKey []byte
	// Public key for verifying access tokens
	AccessTokenPublicKey []byte
	// Maximum age of access tokens in minutes
	AccessTokenMaxAge int
	// Duration for access token expiration
	AccessTokenExpiresIn time.Duration
	// Private key for signing refresh tokens
	RefreshTokenPrivateKey []byte
	// Public key for verifying refresh tokens
	RefreshTokenPublicKey []byte
	// Maximum age of refresh tokens in minutes
	RefreshTokenMaxAge int
	// Duration for refresh token expiration
	RefreshTokenExpiresIn time.Duration
	//
	AllowOrigins string
	//
	ShutdownShelfDetectorChan chan struct{}
}

// New initializes a Config object with values from environment variables and
// certificate files. It returns an error if any of the required environment
// variables or certificate files are missing or inaccessible.
//
// Example:
//
//	cfg, err := config.New()
//	if err != nil {
//	    log.Fatalf("Failed to load configuration: %v", err)
//	}
func New() (*Config, error) {
	certFolder := os.Getenv("CERT_PATH")
	accessPem, err := os.ReadFile(fmt.Sprintf("%s/access.pem", certFolder))
	if err != nil {
		return nil, err
	}
	accessPub, err := os.ReadFile(fmt.Sprintf("%s/access.pub", certFolder))
	if err != nil {
		return nil, err
	}
	refreshPem, err := os.ReadFile(fmt.Sprintf("%s/refresh.pem", certFolder))
	if err != nil {
		return nil, err
	}
	refreshPub, err := os.ReadFile(fmt.Sprintf("%s/refresh.pub", certFolder))
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		API: struct {
			Name string
			Port string
		}{
			Name: os.Getenv("API_NAME"),
			Port: os.Getenv("PORT"),
		},
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
		}{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Cache: struct {
			Host     string
			Port     string
			User     string
			Password string
		}{
			Host:     os.Getenv("CACHE_HOST"),
			Port:     os.Getenv("CACHE_PORT"),
			User:     os.Getenv("CACHE_USER"),
			Password: os.Getenv("CACHE_PASSWORD"),
		},
		AccessTokenPrivateKey:  accessPem,
		AccessTokenPublicKey:   accessPub,
		AccessTokenMaxAge:      15,
		AccessTokenExpiresIn:   time.Minute * 15,
		RefreshTokenPrivateKey: refreshPem,
		RefreshTokenPublicKey:  refreshPub,
		RefreshTokenMaxAge:     60,
		RefreshTokenExpiresIn:  time.Hour * 1,
		AllowOrigins: strings.Join(
			strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
			", ",
		),
		ShutdownShelfDetectorChan: make(chan struct{}, 1),
	}
	return cfg, nil
}
