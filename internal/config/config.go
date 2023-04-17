package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	API struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"api"`
	Database struct {
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Name string `yaml:"name"`
	} `yaml:"db"`
	RSAPrivateKey []byte `yaml:"-"`
	RSAPublicKey  []byte `yaml:"-"`
}

func New(path string) (*Config, error) {
	certFolder := os.Getenv("CERT_PATH")
	prvKey, err := os.ReadFile(certFolder + "id_rsa")
	if err != nil {
		return nil, err
	}
	pubKey, err := os.ReadFile(certFolder + "id_rsa.pub")
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		RSAPrivateKey: prvKey,
		RSAPublicKey:  pubKey,
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
