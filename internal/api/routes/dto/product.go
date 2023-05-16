package dto

type ProductFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,notblank"`
}

type CreateProductDTO struct {
	Name string `json:"name" validate:"required,gte=2,notblank"`
}

type UpdateProductDTO struct {
	Name string `json:"name" validate:"required,gte=2,notblank"`
}

type FindProductDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
