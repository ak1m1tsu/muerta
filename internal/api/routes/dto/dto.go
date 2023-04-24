package dto

type LoginUserPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignUpUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserPayload struct {
	SignUpUserPayload
	Hash string
}

type TokenPayload struct {
	Name  string
	Roles []string
}

type RecipeStepDTO struct {
	ID    int `json:"id" validate:"required, unique"`
	Place int `json:"place" validate:"required, unique"`
}

type CreateRecipeDTO struct {
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Description string          `json:"description,omitempty" validate:"max=200"`
	Steps       []RecipeStepDTO `json:"steps" validate:"required"`
}
