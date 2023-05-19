package dto

type Ingredient struct {
	ProductID int `json:"id_product" validate:"required,gt=0" example:"1"`
	MeasureID int `json:"id_measure" validate:"required,gt=0" example:"1"`
	Quantity  int `json:"quantity"   validate:"required,gt=0" example:"10"`
}

type CreateIngredient struct {
	ProductID int `json:"id_product" validate:"required,gt=0" example:"1"`
	MeasureID int `json:"id_measure" validate:"required,gt=0" example:"1"`
	Quantity  int `json:"quantity"   validate:"required,gt=0" example:"10"`
}

type FindRecipeIngredient struct {
	Product  FindProduct `json:"product"  exmaple:"FindProductDto{ID=1,Name=Томат}"`
	Measure  FindMeasure `json:"measure"  exmaple:"FindMeasureDto{ID=1,Name=Кг}"`
	Quantity int         `json:"quantity" exmaple:"10"`
}

type UpdateIngredient struct {
	ProductID int `json:"id_product" validate:"omitempty,gt=0" example:"1"`
	MeasureID int `json:"id_measure" validate:"omitempty,gt=0" example:"1"`
	Quantity  int `json:"quantity"   validate:"omitempty,gt=0" example:"10"`
}

type DeleteIngredient struct {
	ProductID int `json:"id_product" validate:"required,gt=0" example:"1"`
}
