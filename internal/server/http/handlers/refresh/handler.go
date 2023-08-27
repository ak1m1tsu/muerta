package refresh

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/romankravchuk/muerta/internal/data"
	"github.com/romankravchuk/muerta/internal/server/response"
)

const refreshCookie = "refresh_token"

type Refresher interface {
	Refresh(token string) (*data.TokenDetails, error)
}

func New(log *slog.Logger, refresher Refresher) func(w http.ResponseWriter, r *http.Request) {
	const op = "server.http.handlers.refresh"

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		tokenCookie, err := r.Cookie(refreshCookie)
		if err != nil {
			msg := "cookie with refresh token not found"

			log.Error(msg, slog.String("error", err.Error()))

			response.Error(w, r, http.StatusBadRequest, msg)

			return
		}

		access, err := refresher.Refresh(tokenCookie.Value)
		if err != nil {
			msg := "unable to refresh access token"

			log.Error(msg, slog.String("error", err.Error()))

			response.Error(w, r, http.StatusInternalServerError, msg)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "access_token",
			Value:   access.Token,
			Path:    "/",
			Expires: time.Now().Add(access.ExpiresAt),
		})

		response.OK(w, r, http.StatusOK, render.M{"access_token": access.Token})
	}
}
