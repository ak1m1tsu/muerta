package models

import "time"

type Recipe struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Steps       []Step
}

type Step struct {
	ID    int    `db:"id,id_step"`
	Name  string `db:"name"`
	Place int    `db:"place"`
}
