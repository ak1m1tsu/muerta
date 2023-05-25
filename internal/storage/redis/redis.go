package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/romankravchuk/muerta/internal/pkg/config"
)

type Client interface {
	Get(context.Context, string) *redis.StringCmd
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	Del(context.Context, ...string) *redis.IntCmd
}

func New(cfg *config.Config) (Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port),
		Username: cfg.Cache.User,
		Password: cfg.Cache.Password,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return client, nil
}
