package dto

import "time"

type StorageFilterDTO struct{ *Paging }
type CreateStorageDTO struct {
	Name        string  `json:"name" validate:"required"`
	Temperature float32 `json:"temperature" validate:"required"`
	Humidity    float32 `json:"humidity" validate:"required"`
	TypeID      int     `json:"id_type" validate:"required"`
}
type FindStorageDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Temperature float32    `json:"temperature"`
	Humidity    float32    `json:"humidity"`
	TypeName    string     `json:"type"`
	CreatedAt   *time.Time `json:"created_at"`
}
type UpdateStorageDTO struct {
	Name        string  `json:"name,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	Humidity    float32 `json:"humidity,omitempty"`
	TypeID      int     `json:"id_type,omitempty"`
}
