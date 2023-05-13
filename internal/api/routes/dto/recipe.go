package dto

type CreateRecipeDTO struct {
	UserID      int             `json:"id_user" validate:"required"`
	Name        string          `json:"name" validate:"required,max=100"`
	Description string          `json:"description,omitempty" validate:"max=200"`
	Steps       []RecipeStepDTO `json:"steps" validate:"required"`
}

type FindRecipeDTO struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Steps       []FindStepDTO `json:"steps,omitempty"`
}

type UpdateRecipeDTO struct {
	Name        string `json:"name,omitempty" validate:"max=100"`
	Description string `json:"description,omitempty" validate:"max=200"`
}

type RecipeFilterDTO struct {
	Paging
	Name string `query:"name"`
}

func (f *RecipeFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *RecipeFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *RecipeFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *RecipeFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type FindStepDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Place int    `json:"place,omitempty"`
}

type RecipeStepDTO struct {
	ID    int `json:"id" validate:"required, unique"`
	Place int `json:"place" validate:"required, unique"`
}

type CreateRecipeIngredientDTO struct {
	ProductID int `json:"id_product" validate:"required"`
	MeasureID int `json:"id_measure" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}

type FindRecipeIngredientDTO struct {
	Product  FindProductDTO `json:"product"`
	Measure  FindMeasureDTO `json:"measure"`
	Quantity int            `json:"quantity"`
}

type UpdateRecipeIngredientDTO struct {
	ProductID int `json:"id_product" validate:"required"`
	MeasureID int `json:"id_measure"`
	Quantity  int `json:"quantity" validate:"required"`
}

type DeleteRecipeIngredientDTO struct {
	ProductID int `json:"id_product" validate:"required"`
}

type DeleteRecipeStepDTO struct {
	Place int `json:"place" validate:"required"`
}

type CreateRecipeStepDTO struct {
	Place int `json:"place" validate:"required"`
}
