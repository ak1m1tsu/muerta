package users

import (
	"context"
	"errors"

	"github.com/romankravchuk/muerta/internal/v2/data"
)

var (
	ErrNotFound = errors.New("the user not found")
	ErrExists   = errors.New("the user already exists")
)

type Storage interface {
	Create(ctx context.Context, user *data.User) error
	Update(ctx context.Context, user *data.User) error
	Delete(ctx context.Context, id string) error
	FindByEmail(ctx context.Context, email string) (*data.User, error)
	FindMany(ctx context.Context, filter data.UserFilter) ([]data.User, error)
}
