package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/v2/data"
)

func CreateToken(payload *data.TokenPayload, ttl time.Duration, prvKey []byte) (*data.TokenDetails, error) {
	now := time.Now().UTC()
	td := &data.TokenDetails{
		ID:        uuid.New(),
		Payload:   *payload,
		ExpiresAt: now.Sub(now.Add(ttl)),
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return nil, err
	}
	claims := data.Claims{
		TokenID: td.ID.String(),
		UserID:  payload.ID.String(),
		Email:   payload.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        td.ID.String(),
		},
	}

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, &claims).SignedString(key)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ValidateToken(token string, publickKey []byte) (*data.TokenPayload, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publickKey)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.ParseWithClaims(
		token,
		&data.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return key, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*data.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	payload := &data.TokenPayload{
		ID:     uuid.MustParse(claims.TokenID),
		UserID: uuid.MustParse(claims.UserID),
		Email:  claims.Email,
	}

	return payload, nil
}
