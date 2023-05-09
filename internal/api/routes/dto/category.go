package dto

import "time"

type ProductCategoryFilterDTO struct {
	Paging
	Name string `query:"name"`
}

func (f *ProductCategoryFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *ProductCategoryFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *ProductCategoryFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *ProductCategoryFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type CreateProductCategoryDTO struct {
	Name string `json:"name" validate:"required"`
}

type UpdateProductCategoryDTO struct {
	Name string `json:"name" validate:"required"`
}

type FindProductCategoryDTO struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
