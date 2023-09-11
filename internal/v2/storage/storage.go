package storage

import "errors"

const (
	RollbackFailedMessage = "failed to rollback transaction"
)

var (
	ErrDBPoolIsNnil = errors.New("the database pool is nil")
)
