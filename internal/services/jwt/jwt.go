package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/romankravchuk/muerta/internal/pkg/config"
)

type JWTServicer interface {
	CreateToken(payload interface{}) (string, error)
	ValidateToken(token string) (interface{}, error)
}

type JWTService struct {
	privateKey []byte
	publicKey  []byte
	ttl        time.Duration
}

func New(cfg *config.Config) *JWTService {
	return &JWTService{
		privateKey: cfg.RSAPrivateKey,
		publicKey:  cfg.RSAPublicKey,
		ttl:        time.Hour * 1,
	}
}

func (svc *JWTService) CreateToken(payload interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(svc.privateKey)
	if err != nil {
		return "", fmt.Errorf("create token: parse key: %w", err)
	}

	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["dat"] = payload
	claims["exp"] = now.Add(svc.ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create token: sign token: %w", err)
	}
	return token, nil
}

func (svc *JWTService) ValidateToken(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(svc.publicKey)
	if err != nil {
		return nil, fmt.Errorf("validate token: parse key: %w", err)
	}

	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate token: parse token: %w", err)
	}
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("validate token: invalid claims type")
	}
	return claims["dat"], nil
}
