package dto

type CreateMeasure struct {
	Name string `json:"name" validate:"required,gte=1,notblank" example:"кг"`
}

type UpdateMeasure struct {
	Name string `json:"name" validate:"required,gte=1,notblank" example:"л"`
}

type FindMeasure struct {
	ID   int    `json:"id"   example:"1"`
	Name string `json:"name" example:"кг"`
}
