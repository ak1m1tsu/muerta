package params

import "time"

type CreateProductCategory struct {
	Name string `json:"name" validate:"required,gte=2,notblank" example:"Овощь"`
}

type UpdateProductCategory struct {
	Name string `json:"name" validate:"required,gte=2,notblank" example:"Фрукт"`
}

type FindProductCategory struct {
	ID        int        `json:"id"                   example:"1"`
	Name      string     `json:"name"                 example:"Фрукт"`
	CreatedAt *time.Time `json:"created_at,omitempty" example:"2022-01-01T00:00:00Z"`
}
