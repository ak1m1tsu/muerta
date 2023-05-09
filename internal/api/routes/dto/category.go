package dto

import "time"

type CategoryFilterDTO struct {
	Paging
	Name string `query:"name"`
}

func (f *CategoryFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *CategoryFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *CategoryFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *CategoryFilterDTO) SetOffset(offset int) {
	f.Offset = offset
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
