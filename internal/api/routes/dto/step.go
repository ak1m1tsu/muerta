package dto

type CreateStepDTO struct {
	Name string `json:"name" validate:"required,gt=3,notblank"`
}

type UpdateStepDTO struct {
	Name string `json:"name" validate:"required,gt=3,notblank"`
}

type FindStepDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Place int    `json:"place,omitempty"`
}

type RecipeStepDTO struct {
	ID    int `json:"id" validate:"required,unique,gt=0"`
	Place int `json:"place" validate:"required,unique,gt=0"`
}
