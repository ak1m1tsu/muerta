package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/errors"
)

var (
	errValidateToken    = errors.New("validate token")
	errCreateToken      = errors.New("create token")
	errParseToken       = errors.New("parse token")
	errParseKey         = errors.New("parse key")
	errUnexpectedMethod = errors.New("unexpected method")
	errClaimsType       = errors.New("invalid claims type")
	errSubType          = errors.New("invalid sub type")
	errIdType           = errors.New("invalid id type")
	errNameType         = errors.New("invalid name type")
	errRolesType        = errors.New("invalid roles type")
)

// CreateToken creates a new JWT token with the given payload, TTL, and private key.
// Returns the token details and an error, if any.
func CreateToken(payload *dto.TokenPayload, ttl time.Duration, pirvateKey []byte) (*dto.TokenDetails, error) {
	now := time.Now().UTC()
	td := &dto.TokenDetails{
		UUID: uuid.New().String(),
		User: payload,
	}
	td.ExpiresIn = now.Add(ttl).Unix()

	key, err := jwt.ParseRSAPrivateKeyFromPEM(pirvateKey)
	if err != nil {
		return nil, errCreateToken.With(errParseToken).With(err)
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
		return nil, errCreateToken.With(err)
	}
	return td, nil
}

// ValidateToken validates a JWT token with the given public key.
// Returns the token payload and an error, if any.
func ValidateToken(token string, publickKey []byte) (*dto.TokenPayload, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publickKey)
	if err != nil {
		return nil, errValidateToken.With(errParseKey).With(err)
	}
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errUnexpectedMethod.With(fmt.Errorf("%s", t.Header["alg"]))
		}
		return key, nil
	})
	if err != nil {
		return nil, errValidateToken.With(errParseToken).With(err)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errValidateToken.With(errClaimsType)
	}
	sub, ok := claims["sub"].(map[string]any)
	if !ok {
		return nil, errValidateToken.With(errSubType)
	}
	id, ok := sub["id"].(float64)
	if !ok {
		return nil, errValidateToken.With(errIdType)
	}
	name, ok := sub["name"].(string)
	if !ok {
		return nil, errValidateToken.With(errNameType)
	}
	roles, ok := sub["roles"].([]interface{})
	if !ok {
		return nil, errValidateToken.With(errRolesType)
	}
	paylaod := dto.TokenPayload{
		ID:    int(id),
		Name:  name,
		Roles: roles,
	}
	return &paylaod, nil
}
