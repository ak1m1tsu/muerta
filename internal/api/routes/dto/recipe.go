package dto

type CreateRecipeDTO struct {
	UserID      int             `json:"id_user" validate:"required,gt=0"`
	Name        string          `json:"name" validate:"required,gte=2,lte=100,notblank"`
	Description string          `json:"description,omitempty" validate:"lte=200"`
	Steps       []RecipeStepDTO `json:"steps" validate:"required"`
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

type RecipeFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,notblank"`
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
	ID    int `json:"id" validate:"required,unique,gt=0"`
	Place int `json:"place" validate:"required,unique,gt=0"`
}

type CreateIngredientDTO struct {
	ProductID int `json:"id_product" validate:"required,gt=0"`
	MeasureID int `json:"id_measure" validate:"required,gt=0"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}

type FindRecipeIngredientDTO struct {
	Product  FindProductDTO `json:"product"`
	Measure  FindMeasureDTO `json:"measure"`
	Quantity int            `json:"quantity"`
}

type UpdateIngredientDTO struct {
	ProductID int `json:"id_product" validate:"omitempty,gt=0"`
	MeasureID int `json:"id_measure" validate:"omitempty,gt=0"`
	Quantity  int `json:"quantity" validate:"omitempty,gt=0"`
}

type DeleteIngredientDTO struct {
	ProductID int `json:"id_product" validate:"required,gt=0"`
}

type DeleteRecipeStepDTO struct {
	Place int `json:"place" validate:"required,gt=0"`
}

type CreateRecipeStepDTO struct {
	Place int `json:"place" validate:"required,gt=0"`
}
