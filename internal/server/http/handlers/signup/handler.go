package signup

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/romankravchuk/muerta/internal/data"
	"github.com/romankravchuk/muerta/internal/server/response"
	"github.com/romankravchuk/muerta/internal/storage/users"
)

type SignUper interface {
	SignUp(email, password string) (*data.User, error)
	SendWelcomeMessage(email string) error
}

func New(log *slog.Logger, s SignUper) func(w http.ResponseWriter, r *http.Request) {
	const op = "server.http.handlers.signup"

	type req struct {
		Email           string
		Password        string
		ConfirmPassword string
	}

	type user struct {
		ID        string
		Email     string
		CreatedOn string
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

		user, err := s.SignUp(req.Email, req.Password)
		if errors.Is(err, users.ErrExists) {
			msg := fmt.Sprintf("user with email %s already exists", req.Email)

			log.Error(msg, slog.String("error", err.Error()))

			response.Error(w, r, http.StatusBadRequest, msg)

			return
		}
		if err != nil {
			msg := "unable to sign up"

			log.Error(msg, slog.String("error", err.Error()))

			response.Error(w, r, http.StatusInternalServerError, msg)

			return
		}

		if err := s.SendWelcomeMessage(req.Email); err != nil {
			msg := "unable to send welcome message"

			log.Error(msg, slog.String("error", err.Error()), slog.String("email", req.Email))
		}

		response.OK(w, r, http.StatusOK, render.M{"user": user})
	}
}
