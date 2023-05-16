package dto

import "time"

type CreateShelfLifeDTO struct {
	ProductID    int        `json:"id_product" validate:"required,gt=0"`
	UserID       int        `json:"id_user" validate:"required,gt=0"`
	StorageID    int        `json:"id_storage" validate:"required,gt=0"`
	MeasureID    int        `json:"id_measure" validate:"required,gt=0"`
	Quantity     int        `json:"quantity" validate:"required,gt=0"`
	PurchaseDate *time.Time `json:"purchase_date" validate:"required"`
	EndDate      *time.Time `json:"end_date" validate:"required,gtfield=PurchaseDate"`
}

type FindShelfLifeDTO struct {
	ID           int            `json:"id"`
	Product      FindProductDTO `json:"product"`
	Storage      FindStorageDTO `json:"storage"`
	Measure      FindMeasureDTO `json:"measure"`
	Quantity     int            `json:"quantity"`
	PurchaseDate *time.Time     `json:"purchase_date"`
	EndDate      *time.Time     `json:"end_date"`
}

type UpdateShelfLifeDTO struct {
	ProductID    int        `json:"id_product" validate:"omitempty,gt=0"`
	StorageID    int        `json:"id_storage" validate:"omitempty,gt=0"`
	MeasureID    int        `json:"id_measure" validate:"omitempty,gt=0"`
	Quantity     int        `json:"quantity" validate:"omitempty,gt=0"`
	PurchaseDate *time.Time `json:"purchase_date" validate:"required_with=EndDate,ltfield=EndDate"`
	EndDate      *time.Time `json:"end_date" validate:"required_with=PurchaseDate,gtfield=PurchaseDate"`
}
