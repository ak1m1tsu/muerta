package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/auth"
)

type AuthHandler struct {
	svc           service.AuthServicer
	log           *log.Logger
	accessMaxAge  int
	refreshMaxAge int
}

func New(cfg *config.Config, svc service.AuthServicer, log *log.Logger) *AuthHandler {
	return &AuthHandler{
		svc:           svc,
		log:           log,
		accessMaxAge:  cfg.AccessTokenMaxAge,
		refreshMaxAge: cfg.RefreshTokenMaxAge,
	}
}

// SignUp signs up a new user
//
//	@Summary		Sign up a new user
//	@Description	sign up a new user with the given information
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.SignUp	true	"the sign up information"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/auth/sign-up [post]
func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	payload := new(dto.SignUp)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if payload.Password != payload.PasswordConfirm {
		h.log.ClientError(ctx, fmt.Errorf("passwords do not match"))
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.SignUpUser(ctx.Context(), payload); err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			h.log.ClientError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})

	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Login handles the user login request and returns access and refresh tokens.
//
//	@Summary		Login handles the user login request and returns access and refresh tokens.
//	@Description	Login handles the user login request and returns access and refresh tokens.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			login	body		dto.Login	true	"User credentials"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		401		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	payload := new(dto.Login)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	access, refresh, err := h.svc.LoginUser(ctx.Context(), payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access.Token,
		Path:     "/",
		MaxAge:   h.accessMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh.Token,
		Path:     "/",
		MaxAge:   h.refreshMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   h.accessMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
	})
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data: handlers.Data{
			"access_token":  access.Token,
			"refresh_token": refresh.Token,
		},
	})
}

// RefreshAccessToken refreshes the access token for an authenticated user.
//
//	@Summary		Refresh access token
//	@Description	Refreshes the access token using a refresh token cookie.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handlers.HTTPSuccess{data=handlers.Data{access_token=string,refresh_token=string}}
//	@Failure		403	{object}	handlers.HTTPError
//	@Router			/auth/refresh [post]
//	@Security		Bearer
func (h *AuthHandler) RefreshAccessToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		h.log.ClientError(ctx, fmt.Errorf("refresh token not in cookie"))
		return ctx.Status(http.StatusForbidden).
			JSON(handlers.HTTPError{Error: fiber.ErrForbidden.Error()})
	}
	access, err := h.svc.RefreshAccessToken(refreshToken)
	if err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusForbidden).
			JSON(handlers.HTTPError{Error: fiber.ErrForbidden.Error()})
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access.Token,
		Path:     "/",
		MaxAge:   h.accessMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data:    handlers.Data{"access_token": access.Token},
	})
}

// Logout invalidates the user's access and refresh tokens and logs them out.
//
//	@Summary		Logout user
//	@Description	Invalidates access and refresh tokens, logs out user.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handlers.HTTPSuccess
//	@Router			/auth/logout [post]
//	@Security		Bearer
func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	ctx.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expired,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expired,
	})
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}
