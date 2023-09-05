package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/services/users/proto"
	"github.com/romankravchuk/muerta/internal/v2/storage/users"
	"github.com/romankravchuk/muerta/internal/v2/storage/users/memo"
	"github.com/romankravchuk/nix/logger/sl"
	"github.com/romankravchuk/nix/validator"
	"golang.org/x/crypto/bcrypt"
)

type Option func(*Service) error

func WithUsersStorage(users users.Storage) Option {
	return func(s *Service) error {
		s.users = users
		return nil
	}
}

func WithUsersMemoStorage() Option {
	return func(s *Service) error {
		users := memo.New()
		return WithUsersStorage(users)(s)
	}
}

func WithLogger(log *slog.Logger) Option {
	return func(s *Service) error {
		if log == nil {
			return fmt.Errorf("log is nil")
		}

		s.log = log
		return nil
	}
}

type Service struct {
	users users.Storage

	log *slog.Logger

	proto.UnimplementedUsersServiceServer
}

func New(opts ...Option) (*Service, error) {
	s := &Service{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Service) FindByEmail(ctx context.Context, in *proto.FindByEmailRequest) (*proto.UserResponse, error) {
	if err := validator.Validate(in); err != nil {
		msg := "failed to validate request"

		s.log.Error(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			},
			User: nil,
		}, nil
	}

	user, err := s.users.FindByEmail(ctx, in.GetEmail())
	if errors.Is(err, users.ErrUserNotFound) {
		msg := "user not found"

		s.log.Warn(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusNotFound,
				Error:  msg,
			},
			User: nil,
		}, nil
	}
	if err != nil {
		msg := "failed to find user"

		s.log.Error(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  msg,
			},
			User: nil,
		}, nil
	}

	s.log.Info("found user", slog.String("user", user.Email))

	return &proto.UserResponse{
		Meta: &proto.Response{
			Status: http.StatusOK,
		},
		User: &proto.User{
			Id:                user.ID.String(),
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			EncryptedPassword: user.EncryptedPassword,
			CreatedAt:         user.CreatedAt.String(),
			UpdatedAt:         user.UpdatedAt.String(),
			DeletedAt:         user.DeletedAt.String(),
		},
	}, nil
}

func (s *Service) List(ctx context.Context, in *proto.ListRequest) (*proto.UsersResponse, error) {
	if err := validator.Validate(in); err != nil {
		msg := "failed to validate request"

		s.log.Error(msg, sl.Err(err))

		return &proto.UsersResponse{
			Meta: &proto.Response{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			},
		}, nil
	}

	foundUsers, err := s.users.FindMany(ctx, data.UserFilter{
		Pagination: data.Pagination{
			Limit:  int(in.Limit),
			Offset: int(in.Offset),
		},
		FirstName: in.FirstName,
		LastName:  in.LastName,
	})
	if err != nil {
		if errors.Is(err, users.ErrUsersNotFound) {
			msg := "users not found"

			s.log.Error(msg, sl.Err(err))

			return &proto.UsersResponse{
				Meta: &proto.Response{
					Status: http.StatusNotFound,
					Error:  msg,
				},
			}, nil
		}

		msg := "failed to find users"

		s.log.Error(msg, sl.Err(err))

		return &proto.UsersResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  msg,
			},
		}, nil
	}

	pbUsers := make([]*proto.User, len(foundUsers))
	for i, user := range foundUsers {
		pbUsers[i] = &proto.User{
			Id:                user.ID.String(),
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			EncryptedPassword: user.EncryptedPassword,
			CreatedAt:         user.CreatedAt.String(),
			UpdatedAt:         user.UpdatedAt.String(),
			DeletedAt:         user.DeletedAt.String(),
		}
	}

	return &proto.UsersResponse{
		Meta: &proto.Response{
			Status: http.StatusOK,
		},
		Users: pbUsers,
	}, nil
}

func (s *Service) Create(ctx context.Context, in *proto.CreateRequest) (*proto.UserResponse, error) {
	if err := validator.Validate(in); err != nil {
		msg := "failed to validate request"

		s.log.Error(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			},
		}, nil
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		msg := "failed to generate password hash"

		s.log.Error(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  msg,
			},
		}, nil
	}

	user := &data.User{
		FirstName:         in.GetFirstName(),
		LastName:          in.GetLastName(),
		Email:             in.GetEmail(),
		EncryptedPassword: string(passhash),
	}

	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, users.ErrAlreadyExists) {
			msg := "user already exists"

			s.log.Warn(msg, sl.Err(err))

			return &proto.UserResponse{
				Meta: &proto.Response{
					Status: http.StatusBadRequest,
					Error:  msg,
				},
			}, nil
		}

		msg := "failed to create user"

		s.log.Error(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  msg,
			},
		}, nil
	}

	return &proto.UserResponse{
		Meta: &proto.Response{Status: http.StatusOK},
		User: &proto.User{
			Id:                user.ID.String(),
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			EncryptedPassword: user.EncryptedPassword,
			CreatedAt:         user.CreatedAt.String(),
			UpdatedAt:         user.UpdatedAt.String(),
			DeletedAt:         user.DeletedAt.String(),
		},
	}, nil
}

func (s *Service) Update(ctx context.Context, in *proto.UpdateRequest) (*proto.UserResponse, error) {
	if err := validator.Validate(in); err != nil {
		msg := "invalid request body"

		s.log.Error(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			},
		}, nil
	}

	userID, err := uuid.Parse(in.GetId())
	if err != nil {
		msg := "invalid user id"

		s.log.Error(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusBadRequest,
				Error:  msg,
			},
		}, nil
	}

	user := &data.User{
		ID:        userID,
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
	}

	err = s.users.Update(ctx, user)
	if err != nil {
		msg := "failed to update user"

		s.log.Error(msg, sl.Err(err))

		return &proto.UserResponse{
			Meta: &proto.Response{
				Status: http.StatusInternalServerError,
				Error:  msg,
			},
		}, nil
	}

	return &proto.UserResponse{
		Meta: &proto.Response{
			Status: http.StatusOK,
		},
		User: &proto.User{
			Id:                user.ID.String(),
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			EncryptedPassword: user.EncryptedPassword,
			CreatedAt:         user.CreatedAt.String(),
			UpdatedAt:         user.UpdatedAt.String(),
			DeletedAt:         user.DeletedAt.String(),
		},
	}, nil
}

func (s *Service) Delete(ctx context.Context, in *proto.DeleteRequest) (*proto.Response, error) {
	if err := validator.Validate(in); err != nil {
		msg := "invalid request body"

		s.log.Error(msg, sl.Err(err))

		return &proto.Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	if err := s.users.Delete(ctx, in.GetId()); err != nil {
		msg := "failed to delete user"

		s.log.Error(msg, sl.Err(err))

		return &proto.Response{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	return &proto.Response{
		Status: http.StatusOK,
	}, nil
}
