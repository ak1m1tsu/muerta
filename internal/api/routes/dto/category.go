package dto

import "time"

type ProductCategoryFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,alphaunicode"`
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
	Name string `json:"name" validate:"required,gte=2,alphaunicode"`
}

type UpdateProductCategoryDTO struct {
	Name string `json:"name" validate:"required,gte=2,alphaunicode"`
}

type FindProductCategoryDTO struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
