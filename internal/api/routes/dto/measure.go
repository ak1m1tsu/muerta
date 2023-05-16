package dto

type CreateMeasureDTO struct {
	Name string `json:"name" validate:"required,gte=1,notblank" example:"кг"`
}

type UpdateMeasureDTO struct {
	Name string `json:"name" validate:"required,gte=1,notblank" example:"л"`
}

type FindMeasureDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
