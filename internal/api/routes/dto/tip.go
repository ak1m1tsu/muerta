package dto

type FindTipDTO struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type UpdateTipDTO struct {
	Description string `json:"description" validate:"required"`
}

type CreateTipDTO struct {
	Description string `json:"description" validate:"required"`
}

type TipFilterDTO struct {
	*Paging
	Description string `query:"description"`
}
