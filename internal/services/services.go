package services

import "context"

type Counter interface {
	Count(ctx context.Context) (int, error)
}
