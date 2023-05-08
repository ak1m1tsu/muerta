package dto

type ProductFilterDTO struct {
	*Paging
	Name string `query:"name"`
}

type CreateProductDTO struct {
	Name string `json:"name" validate:"required"`
}

type UpdateProductDTO struct {
	Name string `json:"name"`
}

type FindProductDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
