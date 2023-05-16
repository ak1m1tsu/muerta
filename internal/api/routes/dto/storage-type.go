package dto

type FindStorageTypeDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateStorageTypeDTO struct {
	Name string `json:"name" validate:"required,gte=3,notblank"`
}

type CreateStorageTypeDTO struct {
	Name string `json:"name" validate:"required,gte=3,notblank"`
}
