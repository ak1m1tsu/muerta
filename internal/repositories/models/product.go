package models

import "time"

type Product struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
