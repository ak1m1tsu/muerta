package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	"github.com/romankravchuk/muerta/internal/repositories/user"
	"github.com/romankravchuk/muerta/internal/services/jwt"
	"golang.org/x/crypto/argon2"
)

type AuthServicer interface {
	HashPassword(password string) string
	VerifyPassword(hashedPassword string, candidatePassword string) (bool, error)
	FindUser(ctx context.Context, payload *dto.LoginUserPayload) (models.User, error)
}

type AuthService struct {
	repo        user.UserRepositorer
	svc         jwt.JWTServicer
	iterations  uint32
	memory      uint32
	parallelism uint8
	keyLen      uint32
}

func New(svc jwt.JWTServicer, repo user.UserRepositorer) *AuthService {
	return &AuthService{
		svc:         svc,
		repo:        repo,
		iterations:  3,
		memory:      64 * 1024,
		parallelism: 1,
		keyLen:      32,
	}
}

func (s *AuthService) FindUser(ctx context.Context, payload *dto.LoginUserPayload) (models.User, error) {
	user, err := s.repo.FindByName(ctx, payload.Name)
	if err != nil {
		return models.User{}, fmt.Errorf("find user: %w", err)
	}
	hash := s.HashPassword(payload.Password + user.Salt)
	if _, err := s.repo.FindPassword(ctx, hash); err != nil {
		return models.User{}, fmt.Errorf("find user: %w", err)
	}
	success, err := s.VerifyPassword(hash, payload.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("find user: %w", err)
	}
	if !success {
		return models.User{}, fmt.Errorf("find user: password do not verified")
	}
	return user, nil
}

func (s *AuthService) HashPassword(password string) string {
	salt := []byte(uuid.New().String())
	hash := argon2.IDKey([]byte(password), salt, s.iterations, s.memory, s.parallelism, s.keyLen)

	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$agron2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		s.memory,
		s.iterations,
		s.parallelism,
		b64salt,
		b64hash,
	)
	return encodedHash
}

func (s *AuthService) VerifyPassword(password, encodedHash string) (bool, error) {
	salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, fmt.Errorf("verify password: %w", err)
	}
	otherHash := argon2.IDKey([]byte(password), salt, s.iterations, s.memory, s.parallelism, s.keyLen)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) ([]byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, errors.New("the encoded hash is not in the correct format")
	}
	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, errors.New("incompatible version of argon2")
	}
	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}
	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
	}
	return salt, hash, nil
}
