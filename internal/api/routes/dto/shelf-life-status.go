package dto

type CreateShelfLifeStatusDTO struct {
	Name string `json:"name" validate:"required,gte=3,notblank"`
}

type UpdateShelfLifeStatusDTO struct {
	Name string `json:"name" validate:"required,gte=3,notblank"`
}

type FindShelfLifeStatusDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
