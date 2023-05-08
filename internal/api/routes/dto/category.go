package dto

import "time"

type CategoryFilterDTO struct {
	*Paging
	Name string `query:"name"`
}

type CreateCategoryDTO struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryDTO struct {
	Name string `json:"name" validate:"required"`
}

type FindCategoryDTO struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
}
