package models

import "time"

type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Salt      string    `db:"salt"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
