package dto

import "time"

type ShelfLifeFilterDTO struct {
	Paging
}

func (f *ShelfLifeFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *ShelfLifeFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *ShelfLifeFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *ShelfLifeFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type CreateShelfLifeDTO struct {
	ProductID    int        `json:"id_product" validate:"required"`
	StorageID    int        `json:"id_storage" validate:"required"`
	MeasureID    int        `json:"id_measure" validate:"required"`
	Quantity     int        `json:"quantity" validate:"required"`
	PurchaseDate *time.Time `json:"purchase_date" validate:"required"`
	EndDate      *time.Time `json:"end_date" validate:"required"`
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
	ProductID    int        `json:"id_product"`
	StorageID    int        `json:"id_storage"`
	MeasureID    int        `json:"id_measure"`
	Quantity     int        `json:"quantity"`
	PurchaseDate *time.Time `json:"purchase_date"`
	EndDate      *time.Time `json:"end_date"`
}
