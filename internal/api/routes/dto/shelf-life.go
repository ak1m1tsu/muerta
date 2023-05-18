package dto

import "time"

type CreateShelfLife struct {
	ProductID    int        `json:"id_product"    validate:"required,gt=0"                 example:"1"`
	UserID       int        `json:"id_user"       validate:"required,gt=0"                 example:"1"`
	StorageID    int        `json:"id_storage"    validate:"required,gt=0"                 example:"1"`
	MeasureID    int        `json:"id_measure"    validate:"required,gt=0"                 example:"1"`
	Quantity     int        `json:"quantity"      validate:"required,gt=0"                 example:"1"`
	PurchaseDate *time.Time `json:"purchase_date" validate:"required"                      example:"2020-01-01T00:00:00Z"`
	EndDate      *time.Time `json:"end_date"      validate:"required,gtfield=PurchaseDate"                                exmple:"2020-01-02T00:00:00Z"`
}

type FindShelfLife struct {
	ID           int         `json:"id"            example:"1"`
	Product      FindProduct `json:"product"`
	Storage      FindStorage `json:"storage"`
	Measure      FindMeasure `json:"measure"`
	Quantity     int         `json:"quantity"      example:"1"`
	PurchaseDate *time.Time  `json:"purchase_date" example:"2020-01-01T00:00:00Z"`
	EndDate      *time.Time  `json:"end_date"      example:"2020-01-02T00:00:00Z"`
}

type UpdateShelfLife struct {
	ProductID    int        `json:"id_product"    validate:"omitempty,gt=0"                                  example:"1"`
	StorageID    int        `json:"id_storage"    validate:"omitempty,gt=0"                                  example:"1"`
	MeasureID    int        `json:"id_measure"    validate:"omitempty,gt=0"                                  example:"1"`
	Quantity     int        `json:"quantity"      validate:"omitempty,gt=0"                                  example:"1"`
	PurchaseDate *time.Time `json:"purchase_date" validate:"required_with=EndDate,ltfield=EndDate"           example:"2020-01-01T00:00:00Z"`
	EndDate      *time.Time `json:"end_date"      validate:"required_with=PurchaseDate,gtfield=PurchaseDate" example:"2020-01-02T00:00:00Z"`
}
