package client

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/server/response"
	"github.com/romankravchuk/muerta/internal/v2/services/users/proto"
	"github.com/romankravchuk/nix/logger/sl"
	"github.com/romankravchuk/nix/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client proto.UsersServiceClient
	log    *slog.Logger
}

func New(url string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.Dial(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: proto.NewUsersServiceClient(conn),
		log:    log,
	}, nil
}

func (c *Client) FindByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		msg := "email is required"

		c.log.Error(msg)

		response.Error(w, r, http.StatusBadRequest, msg)

		return
	}

	resp, err := c.client.FindByEmail(r.Context(), &proto.FindByEmailRequest{
		Email: email,
	})
	if err != nil {
		msg := "failed to find user by email"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusInternalServerError, msg)

		return
	}
	if resp.GetMeta().GetError() != "" {
		msg := resp.GetMeta().GetError()

		c.log.Error(msg, slog.String("error", resp.Meta.GetError()))

		response.Error(w, r, int(resp.GetMeta().GetStatus()), msg)

		return
	}

	response.OK(w, r, http.StatusOK, render.M{
		"meta": resp.GetMeta(),
		"user": resp.GetUser(),
	})
}

func (c *Client) List(w http.ResponseWriter, r *http.Request) {
	var (
		filter data.UserFilter
		err    error
	)

	filter.FirstName = r.URL.Query().Get("first_name")
	filter.LastName = r.URL.Query().Get("last_name")
	filter.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		filter.Limit = 10
	}
	filter.Offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		filter.Offset = 0
	}
	filter.Deleted, err = strconv.ParseBool(r.URL.Query().Get("deleted"))
	if err != nil {
		filter.Deleted = false
	}

	resp, err := c.client.List(r.Context(), &proto.ListRequest{
		FirstName: filter.FirstName,
		LastName:  filter.LastName,
		Limit:     int64(filter.Limit),
		Offset:    int64(filter.Offset),
		Deleted:   filter.Deleted,
	})
	if err != nil {
		msg := "failed to list users"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusInternalServerError, msg)

		return
	}
	if resp.GetMeta().GetError() != "" {
		msg := resp.GetMeta().GetError()

		c.log.Error(msg, slog.String("error", resp.Meta.GetError()))

		response.Error(w, r, http.StatusInternalServerError, msg)

		return
	}

	response.OK(w, r, http.StatusOK, render.M{
		"meta":  resp.GetMeta(),
		"users": resp.GetUsers(),
	})
}

func (c *Client) Create(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8"`
		FirstName string `json:"first_name" validate:"omitempty,alphaunicode,gte=2"`
		LastName  string `json:"last_name" validate:"omitempty,alphaunicode,gte=2"`
	}

	var payload *req
	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		msg := "failed to decode request"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusBadRequest, msg)

		return
	}

	if err := validator.Validate(payload); err != nil {
		msg := "failed to validate request"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusBadRequest, err.Error())

		return
	}

	resp, err := c.client.Create(r.Context(), &proto.CreateRequest{
		Email:     payload.Email,
		Password:  payload.Password,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	})
	if err != nil {
		msg := "failed to create user"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusInternalServerError, msg)

		return
	}
	if resp.GetMeta().GetError() != "" {
		msg := resp.GetMeta().GetError()

		c.log.Error(msg, slog.String("error", resp.GetMeta().GetError()))

		response.Error(w, r, int(resp.GetMeta().GetStatus()), msg)

		return
	}

	response.OK(w, r, http.StatusOK, render.M{
		"meta": resp.GetMeta(),
		"user": resp.GetUser(),
	})
}

func (c *Client) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		msg := "id is required"

		c.log.Error(msg)

		response.Error(w, r, http.StatusBadRequest, msg)

		return
	}

	type req struct {
		FirstName string `json:"first_name" validate:"omitempty,alphaunicode,gte=2"`
		LastName  string `json:"last_name" validate:"omitempty,alphaunicode,gte=2"`
	}

	var payload *req

	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		msg := "failed to decode request"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusBadRequest, msg)

		return
	}

	if err := validator.Validate(payload); err != nil {
		msg := "failed to validate request"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusBadRequest, err.Error())

		return
	}

	resp, err := c.client.Update(r.Context(), &proto.UpdateRequest{
		Id:        id,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	})
	if err != nil {
		msg := "failed to update user"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusInternalServerError, msg)

		return
	}
	if resp.GetMeta().GetError() != "" {
		msg := "failed to update user"

		c.log.Error(msg, slog.String("error", resp.GetMeta().GetError()))

		response.Error(w, r, http.StatusInternalServerError, resp.Meta.GetError())

		return
	}

	response.OK(w, r, http.StatusOK, render.M{
		"meta": resp.GetMeta(),
		"user": resp.GetUser(),
	})
}

func (c *Client) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		msg := "id is required"

		c.log.Error(msg, slog.String("error", msg))

		response.Error(w, r, http.StatusBadRequest, msg)

		return
	}

	resp, err := c.client.Delete(r.Context(), &proto.DeleteRequest{
		Id: id,
	})
	if err != nil {
		msg := "failed to delete user"

		c.log.Error(msg, sl.Err(err))

		response.Error(w, r, http.StatusInternalServerError, msg)

		return
	}
	if resp.GetError() != "" {
		msg := resp.GetError()

		c.log.Error(msg, slog.String("error", resp.GetError()))

		response.Error(w, r, int(resp.GetStatus()), msg)

		return
	}

	response.OK	(w, r, http.StatusOK, render.M{
		"meta": resp,
	})
}
