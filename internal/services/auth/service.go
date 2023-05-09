package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/auth"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	"github.com/romankravchuk/muerta/internal/repositories/user"
	"github.com/romankravchuk/muerta/internal/services/jwt"
)

type AuthServicer interface {
	SignUpUser(ctx context.Context, payload *dto.SignUpDTO) error
	LoginUser(ctx context.Context, payload *dto.LoginDTO) (string, error)
}

type AuthService struct {
	repo user.UserRepositorer
	svc  jwt.JWTServicer
}

// LoginUser implements AuthServicer
func (s *AuthService) LoginUser(ctx context.Context, payload *dto.LoginDTO) (string, error) {
	model, err := s.repo.FindByName(ctx, payload.Name)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}
	hash := auth.GenerateHashFromPassword(payload.Password, model.Salt)
	if ok := auth.CompareHashAndPassword(payload.Password, model.Salt, hash); !ok {
		return "", fmt.Errorf("invalid name or password")
	}
	dto := &dto.TokenPayload{
		ID:    model.ID,
		Name:  payload.Name,
		Roles: []interface{}{},
	}
	for _, role := range model.Roles {
		dto.Roles = append(dto.Roles, role.Name)
	}
	token, err := s.svc.CreateToken(dto)
	if err != nil {
		return "", err
	}
	return token, nil
}

// SignUpUser implements AuthServicer
func (s *AuthService) SignUpUser(ctx context.Context, payload *dto.SignUpDTO) error {
	if _, err := s.repo.FindByName(ctx, payload.Name); err == nil {
		return fmt.Errorf("user already exists")
	}
	salt := uuid.New().String()
	hash := auth.GenerateHashFromPassword(payload.Password, salt)
	model := models.User{
		Name: payload.Name,
		Salt: salt,
		Password: models.Password{
			Hash: hash,
		},
	}
	if err := s.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

func New(svc jwt.JWTServicer, repo user.UserRepositorer) AuthServicer {
	return &AuthService{
		svc:  svc,
		repo: repo,
	}
}
