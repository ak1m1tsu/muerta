package config

import (
	"fmt"
	"os"
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
	RSAPrivateKey []byte
	RSAPublicKey  []byte
	TokenMaxAge   int
}

func New() (*Config, error) {
	certFolder := os.Getenv("CERT_PATH")
	prvKey, err := os.ReadFile(fmt.Sprintf("%s/id_rsa.pem", certFolder))
	if err != nil {
		return nil, err
	}
	pubKey, err := os.ReadFile(fmt.Sprintf("%s/id_rsa.pub", certFolder))
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		API: struct {
			Name string
			Port string
		}{
			Name: os.Getenv("NAME"),
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
		RSAPrivateKey: prvKey,
		RSAPublicKey:  pubKey,
		TokenMaxAge:   15,
	}
	return cfg, nil
}
