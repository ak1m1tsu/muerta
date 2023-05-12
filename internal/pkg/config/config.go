package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	API struct {
		Name string
		Port string
	}
	Database struct {
		Host string
		Port string
		User string
		Pass string
		Name string
	}
	AccessTokenPrivateKey  []byte
	AccessTokenPublicKey   []byte
	AccessTokenMaxAge      int
	AccessTokenExpiresIn   time.Duration
	RefreshTokenPrivateKey []byte
	RefreshTokenPublicKey  []byte
	RefreshTokenMaxAge     int
	RefreshTokenExpiresIn  time.Duration
}

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
			Host string
			Port string
			User string
			Pass string
			Name string
		}{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASSWORD"),
			Name: os.Getenv("DB_NAME"),
		},
		AccessTokenPrivateKey:  accessPem,
		AccessTokenPublicKey:   accessPub,
		AccessTokenMaxAge:      15,
		AccessTokenExpiresIn:   time.Minute * 15,
		RefreshTokenPrivateKey: refreshPem,
		RefreshTokenPublicKey:  refreshPub,
		RefreshTokenMaxAge:     60,
		RefreshTokenExpiresIn:  time.Hour * 1,
	}
	return cfg, nil	
}
