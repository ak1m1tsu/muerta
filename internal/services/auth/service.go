package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/auth"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/jwt"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	"github.com/romankravchuk/muerta/internal/repositories/role"
	"github.com/romankravchuk/muerta/internal/repositories/user"
)

type JWTCredential struct {
	PrivateKey []byte
	PublicKey  []byte
	TTL        time.Duration
}

type AuthServicer interface {
	SignUpUser(ctx context.Context, payload *dto.SignUp) error
	LoginUser(
		ctx context.Context,
		payload *dto.Login,
	) (*dto.TokenDetails, *dto.TokenDetails, error)
	RefreshAccessToken(refreshToken string) (*dto.TokenDetails, error)
}

type AuthService struct {
	userRepository user.UserRepositorer
	roleRepository role.RoleRepositorer
	refresh        JWTCredential
	access         JWTCredential
}

// RefreshAccessToken implements AuthServicer
func (s *AuthService) RefreshAccessToken(refreshToken string) (*dto.TokenDetails, error) {
	tokenPayload, err := jwt.ValidateToken(refreshToken, s.refresh.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	access, err := jwt.CreateToken(tokenPayload, s.access.TTL, s.access.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}
	return access, nil
}

// LoginUser implements AuthServicer
func (s *AuthService) LoginUser(
	ctx context.Context,
	payload *dto.Login,
) (*dto.TokenDetails, *dto.TokenDetails, error) {
	model, err := s.userRepository.FindByName(ctx, payload.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %w", err)
	}
	hash := auth.GenerateHashFromPassword(payload.Password, model.Salt)
	if ok := auth.CompareHashAndPassword(payload.Password, model.Salt, hash); !ok {
		return nil, nil, fmt.Errorf("invalid name or password")
	}
	tokenPayload := &dto.TokenPayload{
		UserID:   model.ID,
		Username: payload.Name,
		Roles:    []string{},
	}
	for _, role := range model.Roles {
		tokenPayload.Roles = append(tokenPayload.Roles, role.Name)
	}
	access, err := jwt.CreateToken(tokenPayload, s.access.TTL, s.access.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create access token: %w", err)
	}
	refresh, err := jwt.CreateToken(tokenPayload, s.refresh.TTL, s.refresh.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create refresh token: %w", err)
	}
	return access, refresh, nil
}

// SignUpUser implements AuthServicer
func (s *AuthService) SignUpUser(ctx context.Context, payload *dto.SignUp) error {
	if _, err := s.userRepository.FindByName(ctx, payload.Name); err == nil {
		return fmt.Errorf("user already exists")
	}
	role, err := s.roleRepository.FindByName(ctx, "user")
	if err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}
	salt := uuid.New().String()
	hash := auth.GenerateHashFromPassword(payload.Password, salt)
	model := models.User{
		Name:  payload.Name,
		Salt:  salt,
		Roles: []models.Role{role},
		Password: models.Password{
			Hash: hash,
		},
	}
	if err := s.userRepository.Create(ctx, model); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func New(
	cfg *config.Config,
	repo user.UserRepositorer,
	roleRepository role.RoleRepositorer,
) AuthServicer {
	return &AuthService{
		userRepository: repo,
		roleRepository: roleRepository,
		refresh: JWTCredential{
			PrivateKey: cfg.RefreshTokenPrivateKey,
			PublicKey:  cfg.RefreshTokenPublicKey,
			TTL:        cfg.RefreshTokenExpiresIn,
		},
		access: JWTCredential{
			PrivateKey: cfg.AccessTokenPrivateKey,
			PublicKey:  cfg.AccessTokenPublicKey,
			TTL:        cfg.AccessTokenExpiresIn,
		},
	}
}
