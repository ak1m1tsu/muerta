package dto

import "time"

type CreateProductCategoryDTO struct {
	Name string `json:"name" validate:"required,gte=2,notblank"`
}

type UpdateProductCategoryDTO struct {
	Name string `json:"name" validate:"required,gte=2,notblank"`
}

type FindProductCategoryDTO struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
