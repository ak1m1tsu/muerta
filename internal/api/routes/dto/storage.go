package dto

import "time"

type CreateStorageDTO struct {
	Name        string  `json:"name" validate:"required,gte=3,notblank"`
	Temperature float32 `json:"temperature" validate:"required"`
	Humidity    float32 `json:"humidity" validate:"required"`
	TypeID      int     `json:"id_type" validate:"required,gt=0"`
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
	Name        string  `json:"name,omitempty" validate:"omitempty,gte=3,notblank"`
	Temperature float32 `json:"temperature,omitempty" validate:"omitempty,gte=0"`
	Humidity    float32 `json:"humidity,omitempty" validate:"omitempty,gte=0"`
	TypeID      int     `json:"id_type,omitempty" validate:"omitempty,gte=0"`
}
