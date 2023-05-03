package dto

type CreateRecipeDTO struct {
	Name        string          `json:"name" validate:"required,max=100"`
	Description string          `json:"description,omitempty" validate:"max=200"`
	Steps       []RecipeStepDTO `json:"steps" validate:"required"`
}

type FindRecipeDTO struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Steps       []FindStepDTO `json:"steps"`
}

type UpdateRecipeDTO struct {
	Name        string `json:"name,omitempty" validate:"max=100"`
	Description string `json:"description,omitempty" validate:"max=200"`
}

type RecipeFilterDTO struct {
	*Paging
	Name string `query:"name"`
}

type RecipeStepDTO struct {
	ID    int `json:"id" validate:"required, unique"`
	Place int `json:"place" validate:"required, unique"`
}
