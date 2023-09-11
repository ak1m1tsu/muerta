package auth

import (
	"context"
	"encoding/base64"
	errs "errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/lib/errors"
	"github.com/romankravchuk/muerta/internal/v2/lib/jwt"
	"github.com/romankravchuk/muerta/internal/v2/services/auth/proto"
	"github.com/romankravchuk/muerta/internal/v2/storage"
	"github.com/romankravchuk/muerta/internal/v2/storage/sessions"
	sesmemo "github.com/romankravchuk/muerta/internal/v2/storage/sessions/memo"
	"github.com/romankravchuk/muerta/internal/v2/storage/sessions/redis"
	"github.com/romankravchuk/muerta/internal/v2/storage/users"
	umemo "github.com/romankravchuk/muerta/internal/v2/storage/users/memo"
	"github.com/romankravchuk/muerta/internal/v2/storage/users/postgres"
	"golang.org/x/crypto/bcrypt"
)

type Option func(*Service) error

func WithSessionsStorage(sessions sessions.Storage) Option {
	return func(s *Service) error {
		s.sessions = sessions
		return nil
	}
}

func WithSessionsRedisStorage(url string) Option {
	return func(s *Service) error {
		client, err := storage.NewRedisConnection(url)
		if err != nil {
			return err
		}

		sessions, err := redis.New(client)
		if err != nil {
			return err
		}

		return WithSessionsStorage(sessions)(s)
	}
}

func WithSessionsMemoStorage() Option {
	return func(s *Service) error {
		sessions := sesmemo.New()
		return WithSessionsStorage(sessions)(s)
	}
}

func WithUsersStorage(users users.Storage) Option {
	return func(s *Service) error {
		s.users = users
		return nil
	}
}

func WithUserPostgresStorage(url string) Option {
	return func(s *Service) error {
		conn, err := storage.NewPostgresConnection(url)
		if err != nil {
			return err
		}

		users, err := postgres.New(conn)
		if err != nil {
			return err
		}

		return WithUsersStorage(users)(s)
	}
}

func WithUsersMemoStorage() Option {
	return func(s *Service) error {
		users := umemo.New()
		return WithUsersStorage(users)(s)
	}
}

func WithRefreshCredentials(pvKey, pbKey string, ttl time.Duration) Option {
	return func(s *Service) error {
		pem, pub, err := decodeRSAKeys(pvKey, pbKey)
		if err != nil {
			return err
		}

		s.refreshCreds = data.RSACredentials{
			PrivateKey: pem,
			PublicKey:  pub,
			TTL:        ttl,
		}
		return nil
	}
}

func WithAccessCredentials(pvKey, pbKey string, ttl time.Duration) Option {
	return func(s *Service) error {
		pem, pub, err := decodeRSAKeys(pvKey, pbKey)
		if err != nil {
			return err
		}

		s.accessCreds = data.RSACredentials{
			PrivateKey: pem,
			PublicKey:  pub,
			TTL:        ttl,
		}
		return nil
	}
}

func decodeRSAKeys(private, public string) (pem, pub []byte, err error) {
	pem, err = base64.StdEncoding.DecodeString(private)
	if err != nil {
		return
	}

	pub, err = base64.StdEncoding.DecodeString(public)
	if err != nil {
		return
	}
	return
}

func WithLogger(log *slog.Logger) Option {
	return func(s *Service) error {
		s.log = log
		return nil
	}
}

type Service struct {
	log *slog.Logger

	sessions sessions.Storage
	users    users.Storage

	refreshCreds data.RSACredentials
	accessCreds  data.RSACredentials

	proto.UnsafeAuthServiceServer
}

func New(opts ...Option) (*Service, error) {
	const op = "services.auth.New"

	s := &Service{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, errors.WithOp(op, err)
		}
	}

	return s, nil
}

func (s *Service) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	if in.GetEmail() == "" {
		msg := "email is required"

		s.log.Error(msg, slog.String("email", in.GetEmail()))

		return &proto.LoginResponse{
			Status: http.StatusBadRequest,
			Error:  msg,
		}, nil
	}
	if in.GetPassword() == "" {
		msg := "password is required"

		s.log.Error(msg, slog.String("password", in.GetPassword()))

		return &proto.LoginResponse{
			Status: http.StatusBadRequest,
			Error:  msg,
		}, nil
	}

	user, err := s.users.FindByEmail(ctx, in.GetEmail())
	if err != nil {
		if errs.Is(err, users.ErrUserNotFound) {
			msg := users.ErrUserNotFound.Error()

			s.log.Error(msg, slog.String("error", err.Error()))

			return &proto.LoginResponse{
				Status: http.StatusNotFound,
				Error:  msg,
			}, nil
		}
		msg := "failed to find user by email"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	if bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(in.GetPassword())) != nil {
		msg := "invalid password"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.LoginResponse{
			Status: http.StatusBadRequest,
			Error:  msg,
		}, nil
	}

	details, err := jwt.CreateToken(
		&data.TokenPayload{
			ID:     uuid.New(),
			UserID: user.ID,
			Email:  user.Email,
		},
		s.accessCreds.TTL,
		s.accessCreds.PrivateKey,
	)
	if err != nil {
		msg := "failed to generate access token"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	if err := s.sessions.Set(ctx, details); err != nil {
		msg := "failed to set new session for user"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	return &proto.LoginResponse{
		Status: http.StatusOK,
		Token:  details.Token,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, in *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	if in.GetToken() == "" {
		msg := "token is required"

		s.log.Error(msg, slog.String("token", in.GetToken()))

		return &proto.RefreshResponse{
			Status: http.StatusBadRequest,
			Error:  msg,
		}, nil
	}

	payload, err := jwt.ValidateToken(in.GetToken(), s.refreshCreds.PublicKey)
	if err != nil {
		msg := "failed to validate token"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.RefreshResponse{
			Status: http.StatusBadRequest,
			Error:  msg,
		}, nil
	}

	if _, err := s.sessions.Get(ctx, payload.ID.String()); err != nil {
		msg := "refresh token exired"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.RefreshResponse{
			Status: http.StatusBadRequest,
			Error:  msg,
		}, nil
	}

	accessDetails, err := jwt.CreateToken(payload, s.accessCreds.TTL, s.accessCreds.PrivateKey)
	if err != nil {
		msg := "failed to generate access token"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.RefreshResponse{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	if err := s.sessions.Set(ctx, accessDetails); err != nil {
		msg := "failed to set new session for user"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.RefreshResponse{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	return &proto.RefreshResponse{
		Status: http.StatusOK,
		Token:  accessDetails.Token,
	}, nil
}

func (s *Service) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if in.GetEmail() == "" {
		return &proto.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  "email is required",
		}, nil
	}
	if in.GetPassword() == "" {
		return &proto.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  "password is required",
		}, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return &proto.RegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  "failed to generate hash from password",
		}, nil
	}

	err = s.users.Create(ctx, &data.User{
		Email:             in.GetEmail(),
		EncryptedPassword: string(hash),
	})
	if errs.Is(err, users.ErrAlreadyExists) {
		return &proto.RegisterResponse{
			Status: http.StatusBadRequest,
			Error:  users.ErrAlreadyExists.Error(),
		}, nil
	}
	if err != nil {
		return &proto.RegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  "failed to create user",
		}, nil
	}

	return &proto.RegisterResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Service) Validate(ctx context.Context, in *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	if in.GetToken() == "" {
		msg := "token is required"

		s.log.Error(msg, slog.String("token", in.GetToken()))

		return &proto.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  msg,
		}, nil
	}

	payload, err := jwt.ValidateToken(in.GetToken(), s.accessCreds.PublicKey)
	if err != nil {
		msg := "failed to validate token"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.ValidateResponse{
			Status: http.StatusInternalServerError,
			Error:  msg,
		}, nil
	}

	details, err := s.sessions.Get(ctx, payload.Email)
	if err != nil {
		msg := "access token expired"

		s.log.Error(msg, slog.String("error", err.Error()))

		return &proto.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  msg,
		}, nil
	}

	return &proto.ValidateResponse{
		Status: http.StatusOK,
		UserId: details.Payload.ID.String(),
	}, nil
}
