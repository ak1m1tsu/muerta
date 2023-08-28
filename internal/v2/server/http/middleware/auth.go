package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/romankravchuk/muerta/internal/pkg/jwt"
	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/server/response"
)

var ErrAccessTokenIsEmpty = errors.New("access token is empty")

const (
	authHeader  = "Authorization"
	tokenCookie = "access_token"
	authPrefix  = "Bearer "
)

func Auth(log *slog.Logger, rsaPub string) (func(next http.Handler) http.Handler, error) {
	const op = "server.http.middleware.Auth"

	pubKey, err := base64.StdEncoding.DecodeString(rsaPub)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("auth middleware initialized")

	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("op", op),
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			var token string

			tokenFromHeader := r.Header.Get(authHeader)

			if strings.HasPrefix(tokenFromHeader, authPrefix) {
				token = strings.TrimPrefix(tokenFromHeader, authPrefix)
			} else {
				cookie, err := r.Cookie(tokenCookie)
				if err != nil {
					msg := "cookie with access token not found"

					log.Error(msg, slog.String("error", err.Error()))

					response.Error(w, r, http.StatusUnauthorized, msg)

					return
				}

				token = cookie.Value
			}

			if token == "" {
				msg := "access token not found"

				log.Error(msg, slog.String("error", ErrAccessTokenIsEmpty.Error()))

				response.Error(w, r, http.StatusUnauthorized, msg)

				return
			}

			payload, err := jwt.ValidateToken(token, pubKey)
			if err != nil {
				msg := "access token is invalid"

				log.Error(msg, slog.String("error", err.Error()))

				response.Error(w, r, http.StatusUnauthorized, msg)

				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(
				r.Context(),
				data.ContextKeyUser,
				payload.UserID,
			)))
		}

		return http.HandlerFunc(fn)
	}, nil
}
