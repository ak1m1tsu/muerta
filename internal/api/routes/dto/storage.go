package dto

import "time"

type StorageFilterDTO struct{ Paging }

func (f *StorageFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *StorageFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *StorageFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *StorageFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type CreateStorageDTO struct {
	Name        string  `json:"name" validate:"required"`
	Temperature float32 `json:"temperature" validate:"required"`
	Humidity    float32 `json:"humidity" validate:"required"`
	TypeID      int     `json:"id_type" validate:"required"`
}
type FindStorageDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Temperature float32    `json:"temperature,omitempty"`
	Humidity    float32    `json:"humidity,omitempty"`
	TypeName    string     `json:"type,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}
type UpdateStorageDTO struct {
	Name        string  `json:"name,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	Humidity    float32 `json:"humidity,omitempty"`
	TypeID      int     `json:"id_type,omitempty"`
}
