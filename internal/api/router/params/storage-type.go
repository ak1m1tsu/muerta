package params

type FindStorageType struct {
	ID   int    `json:"id,omitempty"   example:"1"`
	Name string `json:"name,omitempty" example:"Для овощей"`
}

type UpdateStorageType struct {
	Name string `json:"name" validate:"required,gte=3,notblank" example:"Для овощей"`
}

type CreateStorageType struct {
	Name string `json:"name" validate:"required,gte=3,notblank" example:"Для овощей"`
}
