package dto

type FindTipDTO struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type UpdateTipDTO struct {
	Description string `json:"description" validate:"required,gte=3,lte=200"`
}

type CreateTipDTO struct {
	Description string `json:"description" validate:"required,gte=3,lte=200"`
}
