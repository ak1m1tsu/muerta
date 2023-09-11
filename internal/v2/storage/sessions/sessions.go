package sessions

import (
	"context"
	"errors"

	"github.com/romankravchuk/muerta/internal/v2/data"
)

type Storage interface {
	Get(ctx context.Context, key string) (*data.TokenDetails, error)
	Set(ctx context.Context, details *data.TokenDetails) error
}

var ErrTokenNotFound = errors.New("the token not found")
