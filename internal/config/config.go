package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var ErrConfigPath = errors.New("the config path is not specified")

type Config struct {
	Env      string   `yaml:"env"`
	Server   server   `yaml:"server"`
	Postgres postgres `yaml:"postgres"`
	Redis    redis    `yaml:"redis"`
	Access   rsa      `yaml:"access"`
	Refresh  rsa      `yaml:"refresh"`
}

type server struct {
	Addr            string        `yaml:"addr" env:"SERVER_ADDR"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"120s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
}

type postgres struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     string `yaml:"port" env:"POSTGRES_PORT"`
	Database string `yaml:"name" env:"POSTGRES_DB"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

type redis struct {
	Host         string        `yaml:"host" env:"REDIS_HOST"`
	Port         string        `yaml:"port" env:"REDIS_PORT"`
	Database     string        `yaml:"name" env:"REDIS_DB"`
	User         string        `yaml:"user" env:"REDIS_USER"`
	Password     string        `yaml:"password" env:"REDIS_PASSWORD"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
	DialTimeout  time.Duration `yaml:"dial_timeout" env-default:"5s"`
}

type rsa struct {
	PublicKey  string        `yaml:"public_key" env:"RSA_PUBLIC_KEY"`
	PrivateKey string        `yaml:"private_key" env:"RSA_PRIVATE_KEY"`
	Expiration time.Duration `yaml:"expiration"`
}

func Load() (*Config, error) {
	const op = "config.Load"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrConfigPath)
	}

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var cfg *Config
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cfg, nil
}
