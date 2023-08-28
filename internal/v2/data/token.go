package data

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenPayload struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Email  string
}

type TokenDetails struct {
	Token     string
	ID        uuid.UUID
	Payload   TokenPayload
	ExpiresAt time.Duration
}

type Claims struct {
	TokenID string `json:"token_id"`
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

type RSACredentials struct {
	PrivateKey []byte
	PublicKey  []byte
	TTL        time.Duration
}
