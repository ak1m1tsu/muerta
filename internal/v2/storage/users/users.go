package users

import (
	"context"
	"errors"

	"github.com/romankravchuk/muerta/internal/v2/data"
)

var (
	ErrUserNotFound  = errors.New("the user not found")
	ErrAlreadyExists = errors.New("the user already exists")
	ErrUsersNotFound = errors.New("the users not found")
)

type Storage interface {
	Create(ctx context.Context, user *data.User) error
	Update(ctx context.Context, user *data.User) error
	Delete(ctx context.Context, id string) error
	FindByEmail(ctx context.Context, email string) (*data.User, error)
	FindMany(ctx context.Context, filter data.UserFilter) ([]data.User, error)
}
