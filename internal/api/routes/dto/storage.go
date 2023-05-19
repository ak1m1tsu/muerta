package dto

import "time"

type CreateStorage struct {
	Name        string  `json:"name"        validate:"required,gte=3,notblank" example:"Холодильник"`
	Temperature float32 `json:"temperature" validate:"required"                example:"1"`
	Humidity    float32 `json:"humidity"    validate:"required"                example:"1"`
	TypeID      int     `json:"id_type"     validate:"required,gt=0"           example:"1"`
}

type FindStorage struct {
	ID          int             `json:"id"                    example:"1"`
	Name        string          `json:"name"                  example:"Холодильник"`
	Temperature float32         `json:"temperature,omitempty" example:"1"`
	Humidity    float32         `json:"humidity,omitempty"    example:"1"`
	Type        FindStorageType `json:"type,omitempty"        example:"1"`
	CreatedAt   *time.Time      `json:"created_at,omitempty"  example:"2022-01-01T00:00:00Z"`
}

type UpdateStorage struct {
	Name        string  `json:"name,omitempty"        validate:"omitempty,gte=3,notblank" example:"Холодильник"`
	Temperature float32 `json:"temperature,omitempty" validate:"omitempty,gte=0"          example:"1"`
	Humidity    float32 `json:"humidity,omitempty"    validate:"omitempty,gte=0"          example:"1"`
	TypeID      int     `json:"id_type,omitempty"     validate:"omitempty,gte=0"          example:"1"`
}
