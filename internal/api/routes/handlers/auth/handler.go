package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	dto.SignUpDTO	true	"the sign up information"
//	@Description	sign up a new user with the given information
//	@Success		200	{object}	handlers.apiResponse
//	@Failure		400	{object}	handlers.errorResponse
//	@Failure		502	{object}	handlers.errorResponse
//	@Router			/auth/sign-up [post]
func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	var payload *dto.SignUpDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	if payload.Password != payload.PasswordConfirm {
		h.log.ClientError(ctx, fmt.Errorf("passwords do not match"))
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	if err := h.svc.SignUpUser(ctx.Context(), payload); err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			h.log.ClientError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
		}
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse())
}

// Login handles the user login request and returns access and refresh tokens.
//
//	@Summary		Login handles the user login request and returns access and refresh tokens.
//	@Description	Login handles the user login request and returns access and refresh tokens.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			login	body		dto.LoginDTO	true	"User credentials"
//	@Success		200		{object}	handlers.apiResponse
//	@Failure		401		{object}	handlers.errorResponse
//	@Failure		502		{object}	handlers.errorResponse
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var payload *dto.LoginDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusUnauthorized).
			JSON(handlers.ErrorResponse(fiber.ErrUnauthorized))
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusUnauthorized).
			JSON(handlers.ErrorResponse(fiber.ErrUnauthorized))
	}
	access, refresh, err := h.svc.LoginUser(ctx.Context(), payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
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
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{
			"access_token":  access.Token,
			"refresh_token": refresh.Token,
		},
	))
}

// RefreshAccessToken refreshes the access token for an authenticated user.
//
//	@Summary		Refresh access token
//	@Description	Refreshes the access token using a refresh token cookie.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handlers.apiResponse
//	@Failure		403	{object}	handlers.errorResponse
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshAccessToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		h.log.ClientError(ctx, fmt.Errorf("refresh token not in cookie"))
		return ctx.Status(http.StatusForbidden).
			JSON(handlers.ErrorResponse(fiber.ErrForbidden))
	}
	access, err := h.svc.RefreshAccessToken(refreshToken)
	if err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusForbidden).
			JSON(handlers.ErrorResponse(fiber.ErrForbidden))
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access.Token,
		Path:     "/",
		MaxAge:   h.accessMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"access_token": access.Token},
	))
}

// Logout invalidates the user's access and refresh tokens and logs them out.
//
//	@Summary		Logout user
//	@Description	Invalidates access and refresh tokens, logs out user.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handlers.apiResponse
//	@Router			/auth/logout [post]
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
	return ctx.JSON(handlers.SuccessResponse())
}
