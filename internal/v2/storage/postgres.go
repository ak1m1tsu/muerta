package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/romankravchuk/muerta/internal/v2/lib/errors"
)

const UniqueViolation = "23505"

func NewPostgresConnection(url string) (*sql.DB, error) {
	const op = "storage.NewPostgresConnection"

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.WithOp(op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, errors.WithOp(op, err)
	}

	return db, nil
}
