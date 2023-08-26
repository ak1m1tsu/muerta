package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const UniqueViolation = "23505"

func NewPostgresConnection(url string) (*sql.DB, error) {
	const op = "storage.NewPostgresConnection"

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
