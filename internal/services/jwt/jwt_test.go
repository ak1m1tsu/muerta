package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

type testPayload struct {
	Name string
	Role string
}

func generateCfg() config.Config {
	rng := rand.Reader
	pk, _ := rsa.GenerateKey(rng, 4096)
	cfg := config.Config{
		RSAPrivateKey: pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(pk),
			},
		),
		RSAPublicKey: pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(&pk.PublicKey),
			},
		),
	}
	return cfg
}

func Test_CreateToken(t *testing.T) {
	cfg := generateCfg()
	svc := New(&cfg)
	payload := testPayload{
		Name: "test_user",
		Role: "user",
	}
	token, err := svc.CreateToken(payload)
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token)
}
