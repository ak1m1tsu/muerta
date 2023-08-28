package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/romankravchuk/muerta/internal/v2/data"
)

var ErrRedisClientIsNil = errors.New("the redis client is nil")

type SessionsStorage struct {
	client *redis.Client
}

func New(client *redis.Client) (*SessionsStorage, error) {
	const op = "storage.sessions.redis.New"

	if client == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrRedisClientIsNil)
	}

	return &SessionsStorage{
		client: client,
	}, nil
}

func (s *SessionsStorage) Get(ctx context.Context, key string) (*data.TokenDetails, error) {
	const op = "storage.sessions.redis.SessionsStorage.Get"

	ctx, cancel := context.WithTimeout(
		ctx,
		s.client.Options().ReadTimeout,
	)
	defer cancel()

	res, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	details := data.TokenDetails{}
	if err := sonic.UnmarshalString(res, &details); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &details, nil
}

func (s *SessionsStorage) Set(ctx context.Context, details *data.TokenDetails) error {
	const op = "storage.sessions.redis.SessionsStorage.Set"

	ctx, cancel := context.WithTimeout(
		ctx,
		s.client.Options().WriteTimeout,
	)
	defer cancel()

	marshaledDetails, err := sonic.MarshalString(details)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.client.Set(
		ctx,
		details.Payload.Email,
		marshaledDetails,
		details.ExpiresAt,
	).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
