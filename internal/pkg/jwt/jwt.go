package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
)

func CreateToken(payload *dto.TokenPayload, ttl time.Duration, pirvateKey []byte) (*dto.TokenDetails, error) {
	now := time.Now().UTC()
	td := &dto.TokenDetails{
		UUID: uuid.New().String(),
		User: payload,
	}
	td.ExpiresIn = now.Add(ttl).Unix()

	key, err := jwt.ParseRSAPrivateKeyFromPEM(pirvateKey)
	if err != nil {
		return nil, fmt.Errorf("create token: parse key: %w", err)
	}

	claims := make(jwt.MapClaims)
	claims["sub"] = map[string]any{
		"id":    payload.ID,
		"name":  payload.Name,
		"roles": payload.Roles,
	}
	claims["token_uuid"] = td.UUID
	claims["exp"] = td.ExpiresIn
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create token: sign token: %w", err)
	}
	return td, nil
}

func ValidateToken(token string, publickKey []byte) (*dto.TokenPayload, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publickKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate token: parse token: %w", err)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("validate token: invalid claims type")
	}
	paylaod := dto.TokenPayload{
		ID:    int(claims["sub"].(map[string]any)["id"].(float64)),
		Name:  claims["sub"].(map[string]any)["name"].(string),
		Roles: claims["sub"].(map[string]any)["roles"].([]interface{}),
	}
	return &paylaod, nil
}
