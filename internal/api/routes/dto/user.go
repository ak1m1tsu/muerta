package dto

import "time"

type FindUserDTO struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	CreatedAt time.Time        `json:"created_at,omitempty"`
	Settings  []FindSettingDTO `json:"settings,omitempty"`
}

type UpdateUserDTO struct {
	Name string `json:"name" validate:"required,gt=3,alpha"`
}

type CreateUserDTO struct {
	Name     string           `json:"name" validate:"required,gte=5,alphanum"`
	Password string           `json:"password" validate:"required,gte=8,alphanum"`
	Settings []UserSettingDTO `json:"settings"`
	Roles    []UserRoleDTO    `json:"roles"`
}

type UserSettingDTO struct {
	ID    int    `json:"id" validate:"required,gt=0"`
	Value string `json:"value" validate:"required,gt=0,notblank"`
}

type UserRoleDTO struct {
	ID int `json:"id"`
}

type UpdateUserSettingDTO struct {
	Value string `json:"value" validate:"required,gt=0,notblank"`
}

type UserStorageDTO struct {
	StorageID int `json:"id_storage" validate:"gt=0"`
}

type UserShelfLifeDTO struct {
	ShelfLifeID  int        `json:"id_shelf_life" validate:"gt=0"`
	ProductID    int        `json:"id_product" validate:"omitempty,gt=0"`
	StorageID    int        `json:"id_storage" validate:"omitempty,gt=0"`
	MeasureID    int        `json:"id_measure" validate:"omitempty,gt=0"`
	Quantity     int        `json:"quantity" validate:"omitempty,gt=0"`
	PurchaseDate *time.Time `json:"purchase_date" validate:"required_with=EndDate,ltfield=EndDate"`
	EndDate      *time.Time `json:"end_date" validate:"required_with=PurchaseDate,gtfield=PurchaseDate"`
}
