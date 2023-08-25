package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(url string) (*redis.Client, error) {
	const op = "storage.NewRedisConnection"

	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := redis.NewClient(options)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return client, err
}
