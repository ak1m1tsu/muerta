package models

import "time"

type User struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Salt      string `db:"salt"`
	Password  Password
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Password struct {
	Hash string `db:"password"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
