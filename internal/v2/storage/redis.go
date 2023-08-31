package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/romankravchuk/muerta/internal/v2/lib/errors"
)

func NewRedisConnection(url string) (*redis.Client, error) {
	const op = "storage.NewRedisConnection"

	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.WithOp(op, err)
	}

	client := redis.NewClient(options)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.WithOp(op, err)
	}

	return client, err
}
