package dto

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
