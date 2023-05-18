package dto

type FindStorageType struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Для овощей"`
}

type UpdateStorageType struct {
	Name string `json:"name" validate:"required,gte=3,notblank" example:"Для овощей"`
}

type CreateStorageType struct {
	Name string `json:"name" validate:"required,gte=3,notblank" example:"Для овощей"`
}
