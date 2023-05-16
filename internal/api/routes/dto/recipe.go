package dto

type CreateRecipeDTO struct {
	UserID      int             `json:"id_user" validate:"required,gt=0"`
	Name        string          `json:"name" validate:"required,gte=2,lte=100,notblank"`
	Description string          `json:"description,omitempty" validate:"lte=200"`
	Steps       []RecipeStepDTO `json:"steps" validate:"required"`
	Ingredients []IngredientDTO `json:"ingredients" validate:"required"`
}

type FindRecipeDTO struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Steps       []FindStepDTO `json:"steps,omitempty"`
}

type UpdateRecipeDTO struct {
	Name        string `json:"name" validate:"gte=2,lte=100"`
	Description string `json:"description" validate:"lte=200"`
}

type DeleteRecipeStepDTO struct {
	Place int `json:"place" validate:"required,gt=0"`
}

type CreateRecipeStepDTO struct {
	Place int `json:"place" validate:"required,gt=0"`
}
