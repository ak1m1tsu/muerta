package params

type CreateShelfLifeStatus struct {
	Name string `json:"name" validate:"required,gte=3,notblank" example:"Просрочен"`
}

type UpdateShelfLifeStatus struct {
	Name string `json:"name" validate:"required,gte=3,notblank" example:"Просрочен"`
}

type FindShelfLifeStatus struct {
	ID   int    `json:"id"   example:"1"`
	Name string `json:"name" example:"Просрочен"`
}
