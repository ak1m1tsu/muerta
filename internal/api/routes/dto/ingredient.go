package dto

type IngredientDTO struct {
	ProductID int `json:"id_product" validate:"required,gt=0"`
	MeasureID int `json:"id_measure" validate:"required,gt=0"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
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
