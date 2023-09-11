package redis

import (
	"context"
	errs "errors"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/lib/errors"
)

var ErrRedisClientIsNil = errs.New("the redis client is nil")

type Storage struct {
	client *redis.Client
}

func New(client *redis.Client) (*Storage, error) {
	const op = "storage.sessions.redis.New"

	if client == nil {
		return nil, errors.WithOp(op, ErrRedisClientIsNil)
	}

	return &Storage{
		client: client,
	}, nil
}

func (s *Storage) Get(ctx context.Context, key string) (*data.TokenDetails, error) {
	const op = "storage.sessions.redis.SessionsStorage.Get"

	ctx, cancel := context.WithTimeout(
		ctx,
		s.client.Options().ReadTimeout,
	)
	defer cancel()

	res, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.WithOp(op, err)
	}

	details := data.TokenDetails{}
	if err := sonic.UnmarshalString(res, &details); err != nil {
		return nil, errors.WithOp(op, err)
	}

	return &details, nil
}

func (s *Storage) Set(ctx context.Context, details *data.TokenDetails) error {
	const op = "storage.sessions.redis.SessionsStorage.Set"

	ctx, cancel := context.WithTimeout(
		ctx,
		s.client.Options().WriteTimeout,
	)
	defer cancel()

	marshaledDetails, err := sonic.MarshalString(details)
	if err != nil {
		return errors.WithOp(op, err)
	}

	err = s.client.Set(
		ctx,
		details.Payload.Email,
		marshaledDetails,
		details.ExpiresAt,
	).Err()
	if err != nil {
		return errors.WithOp(op, err)
	}

	return nil
}
