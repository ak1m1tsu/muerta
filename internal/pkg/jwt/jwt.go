package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/pkg/errors"
)

var (
	errValidateToken    = errors.New("validate token")
	errCreateToken      = errors.New("create token")
	errParseToken       = errors.New("parse token")
	errParseKey         = errors.New("parse key")
	errUnexpectedMethod = errors.New("unexpected method")
	errClaimsType       = errors.New("invalid claims type")
)

type Claims struct {
	UserID   int      `json:"user_id,omitempty"`
	Username string   `json:"username,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// CreateToken creates a new JWT token with the given payload, TTL, and private key.
// Returns the token details and an error, if any.
func CreateToken(
	payload *params.TokenPayload,
	ttl time.Duration,
	pirvateKey []byte,
) (*params.TokenDetails, error) {
	now := time.Now().UTC()
	td := &params.TokenDetails{
		UUID: uuid.New().String(),
		User: payload,
	}
	td.ExpiresIn = now.Add(ttl).Unix()

	key, err := jwt.ParseRSAPrivateKeyFromPEM(pirvateKey)
	if err != nil {
		return nil, errCreateToken.With(errParseToken).With(err)
	}
	claims := Claims{
		UserID:   payload.UserID,
		Username: payload.Username,
		Roles:    payload.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        td.UUID,
		},
	}

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, &claims).SignedString(key)
	if err != nil {
		return nil, errCreateToken.With(err)
	}
	return td, nil
}

// ValidateToken validates a JWT token with the given public key.
// Returns the token payload and an error, if any.
func ValidateToken(token string, publickKey []byte) (*params.TokenPayload, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publickKey)
	if err != nil {
		return nil, errValidateToken.With(errParseKey).With(err)
	}
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errUnexpectedMethod.With(fmt.Errorf("%s", t.Header["alg"]))
			}
			return key, nil
		},
	)
	if err != nil {
		return nil, errValidateToken.With(errParseToken).With(err)
	}
	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return nil, errValidateToken.With(errClaimsType)
	}
	payload := &params.TokenPayload{
		UUID:     claims.ID,
		UserID:   claims.UserID,
		Username: claims.Username,
		Roles:    claims.Roles,
	}
	return payload, nil
}
