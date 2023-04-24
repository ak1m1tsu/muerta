package models

import "time"

type Recipe struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Steps       []RecipeStep
}

type RecipeStep struct {
	RecipeID int `db:"id_recipe"`
	StepID   int `db:"id_step"`
	Place    int `db:"place"`
}

type Step struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
