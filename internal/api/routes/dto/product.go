package dto

type CreateProduct struct {
	Name string `json:"name" validate:"required,gte=2,notblank" example:"Томат"`
}

type UpdateProduct struct {
	Name string `json:"name" validate:"required,gte=2,notblank" exmaple:"Морковь"`
}

type FindProduct struct {
	ID   int    `json:"id"   example:"1"`
	Name string `json:"name" example:"Морковь"`
}
