package token

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/server/response"
	"github.com/romankravchuk/muerta/internal/v2/storage/users"
)

const (
	accessCookie  = "access_token"
	refreshCookie = "refresh_token"
)

type Signiner interface {
	SignIn(email, password string) (access *data.TokenDetails, refresh *data.TokenDetails, err error)
}

func New(log *slog.Logger, signiner Signiner) func(w http.ResponseWriter, r *http.Request) {
	const op = "server.http.handlers.token"

	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req *req
		if err := render.DecodeJSON(r.Body, req); err != nil {
			msg := "unable to decode request"

			log.Error(msg, slog.String("error", err.Error()))

			response.Error(w, r, http.StatusBadRequest, msg)

			return
		}

		// if err := validator.Validate(req); err != nil {
		// 	msg := "failed to validate request"

		// 	log.Error(msg, slog.String("error", err.Error()))

		// 	response.Error(w, r, http.StatusBadRequest, err.Error())

		// 	return
		// }

		access, refresh, err := signiner.SignIn(req.Email, req.Password)
		if errors.Is(err, users.ErrNotFound) {
			msg := fmt.Sprintf("user with email %s not found", req.Email)

			log.Error(msg, slog.String("error", err.Error()))

			response.Error(w, r, http.StatusNotFound, msg)

			return
		} else if err != nil {
			msg := "unable to sign in"

			log.Error(msg, slog.String("error", err.Error()))

			response.Error(w, r, http.StatusInternalServerError, msg)

			return
		}

		now := time.Now()
		http.SetCookie(w, &http.Cookie{
			Name:    accessCookie,
			Value:   access.Token,
			Path:    "/",
			Expires: now.Add(access.ExpiresAt),
		})

		http.SetCookie(w, &http.Cookie{
			Name:    refreshCookie,
			Value:   refresh.Token,
			Path:    "/",
			Expires: now.Add(refresh.ExpiresAt),
		})

		response.OK(w, r, http.StatusOK, render.M{
			accessCookie:  access.Token,
			refreshCookie: refresh.Token,
		})
	}
}
