package utils

import (
	"errors"
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

const AuthorizationHeader = "Authorization"

const AuthorizationCookie = "access_token"

var ErrNoToken = errors.New("the access token not found")

func GetTokenFromReq(r *http.Request) (string, error) {
	var token string

	auth := r.Header.Get(AuthorizationHeader)
	if strings.HasPrefix(auth, bearerPrefix) {
		token = strings.TrimPrefix(auth, bearerPrefix)
	} else {
		cookie, err := r.Cookie(AuthorizationCookie)
		if err != nil {
			return "", err
		}

		token = cookie.Value
	}

	if token == "" {
		return "", ErrNoToken
	}

	return token, nil
}
