package data

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `db:"id"`
	FirstName         string    `db:"first_name"`
	LastName          string    `db:"last_name"`
	Email             string    `db:"email"`
	EncryptedPassword string    `db:"encrypted_password"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	DeletedAt         time.Time `db:"deleted_at"`
}

type UserFilter struct {
	Pagination
	FirstName string
	LastName  string
}
