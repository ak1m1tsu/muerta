package params

type FindTip struct {
	ID          int    `json:"id"          example:"1"`
	Description string `json:"description" example:"Хранить в Холодильнике"`
}

type UpdateTip struct {
	Description string `json:"description" validate:"required,gte=3,lte=200" example:"Хранить в Холодильнике"`
}

type CreateTip struct {
	Description string `json:"description" validate:"required,gte=3,lte=200" example:"Хранить в Холодильнике"`
}
