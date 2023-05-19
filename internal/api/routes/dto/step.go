package dto

type CreateStep struct {
	Name string `json:"name" validate:"required,gt=3,notblank" example:"Сварить картошку"`
}

type UpdateStep struct {
	Name string `json:"name" validate:"required,gt=3,notblank" example:"Сварить картошку"`
}

type FindStep struct {
	ID    int    `json:"id"              example:"1"`
	Name  string `json:"name"            example:"Сварить картошку"`
	Place int    `json:"place,omitempty" example:"1"`
}

type RecipeStep struct {
	ID    int `json:"id"    validate:"required,unique,gt=0" example:"1"`
	Place int `json:"place" validate:"required,unique,gt=0" example:"1"`
}
