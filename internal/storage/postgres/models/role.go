package models

type Role struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
