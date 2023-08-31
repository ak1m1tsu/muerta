package client

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/romankravchuk/muerta/internal/v2/lib/utils"
	"github.com/romankravchuk/muerta/internal/v2/server/response"
	"github.com/romankravchuk/muerta/internal/v2/services/auth/proto"
	"github.com/romankravchuk/nix/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client proto.AuthServiceClient
	log    *slog.Logger
}

func New(url string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: proto.NewAuthServiceClient(conn),
		log:    log,
	}, nil
}

func (c *Client) Register(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,gte=8,alphanum"`
		PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
	}

	var payload *req
	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())

		c.log.Error("failed to decode request body", slog.String("error", err.Error()))

		return
	}

	if err := validator.Validate(payload); err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())

		c.log.Error("failed to validate request body", slog.String("error", err.Error()))

		return
	}

	resp, err := c.client.Register(r.Context(), &proto.RegisterRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())

		c.log.Error("failed to register new user", slog.String("error", err.Error()))

		return
	}
	if resp.GetError() != "" {
		response.Error(w, r, int(resp.GetStatus()), resp.GetError())

		c.log.Error("failed to register new user", slog.String("error", resp.GetError()))

		return
	}

	response.OK(w, r, int(resp.GetStatus()), render.M{
		"status": resp.GetStatus(),
	})
}

func (c *Client) Login(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	var payload *req
	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())

		c.log.Error("failed to decode request body", slog.String("error", err.Error()))

		return
	}

	if err := validator.Validate(payload); err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())

		c.log.Error("failed to validate request body", slog.String("error", err.Error()))

		return
	}

	resp, err := c.client.Login(r.Context(), &proto.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())

		c.log.Error("failed to make request to auth service", slog.String("error", err.Error()))

		return
	}
	if resp.GetError() != "" {
		response.Error(w, r, int(resp.GetStatus()), resp.GetError())

		c.log.Error("failed to make request to auth service", slog.String("error", resp.GetError()))

		return
	}

	response.OK(w, r, int(resp.GetStatus()), render.M{
		"token": resp.GetToken(),
	})
}

func (c *Client) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetTokenFromReq(r)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, err.Error())

		c.log.Error("failed to get token from request", slog.String("error", err.Error()))

		return
	}

	resp, err := c.client.Refresh(r.Context(), &proto.RefreshRequest{
		Token: token,
	})
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())

		c.log.Error("failed to make request to auth service", slog.String("error", err.Error()))

		return
	}
	if resp.GetError() != "" {
		response.Error(w, r, int(resp.GetStatus()), resp.GetError())

		c.log.Error("failed to make request to auth service", slog.String("error", resp.GetError()))

		return
	}

	response.OK(w, r, int(resp.GetStatus()), render.M{
		"status": resp.GetStatus(),
		"token":  resp.GetToken(),
	})
}

func (c *Client) Validate(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetTokenFromReq(r)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, err.Error())

		c.log.Error("failed to get token from request", slog.String("error", err.Error()))

		return
	}

	resp, err := c.client.Validate(r.Context(), &proto.ValidateRequest{
		Token: token,
	})
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())

		c.log.Error("failed to make request to auth service", slog.String("error", err.Error()))

		return
	}
	if resp.GetError() != "" {
		response.Error(w, r, int(resp.GetStatus()), resp.GetError())

		c.log.Error("failed to make request to auth service", slog.String("error", resp.GetError()))

		return
	}

	response.OK(w, r, http.StatusOK, render.M{
		"user_id": resp.GetUserId(),
	})
}
