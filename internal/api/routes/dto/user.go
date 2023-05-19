package dto

import "time"

type FindUser struct {
	ID        int           `json:"id"                   example:"1"`
	Name      string        `json:"name"                 example:"user"`
	CreatedAt time.Time     `json:"created_at,omitempty" example:"2020-01-01T00:00:00Z"`
	Settings  []FindSetting `json:"settings,omitempty"`
}

type UpdateUser struct {
	Name string `json:"name" validate:"required,gt=3,alpha" example:"user"`
}

type CreateUser struct {
	Name     string        `json:"name"     validate:"required,gte=5,alphanum" example:"user"`
	Password string        `json:"password" validate:"required,gte=8,alphanum" example:"th3B3stUs3r"`
	Settings []UserSetting `json:"settings"`
	Roles    []UserRole    `json:"roles"`
}

type UserSetting struct {
	ID    int    `json:"id"    validate:"required,gt=0"          example:"1"`
	Value string `json:"value" validate:"required,gt=0,notblank" example:"Да"`
}

type UserRole struct {
	ID int `json:"id" example:"1"`
}

type UpdateUserSetting struct {
	Value string `json:"value" validate:"required,gt=0,notblank" example:"Нет"`
}

type UserStorage struct {
	StorageID int `json:"id_storage" validate:"gt=0" example:"1"`
}

type UserShelfLife struct {
	ShelfLifeID  int        `json:"id_shelf_life" validate:"gt=0"                                            example:"1"`
	ProductID    int        `json:"id_product"    validate:"omitempty,gt=0"                                  example:"1"`
	StorageID    int        `json:"id_storage"    validate:"omitempty,gt=0"                                  example:"1"`
	MeasureID    int        `json:"id_measure"    validate:"omitempty,gt=0"                                  example:"1"`
	Quantity     int        `json:"quantity"      validate:"omitempty,gt=0"                                  example:"1"`
	PurchaseDate *time.Time `json:"purchase_date" validate:"required_with=EndDate,ltfield=EndDate"           example:"2020-01-01T00:00:00Z"`
	EndDate      *time.Time `json:"end_date"      validate:"required_with=PurchaseDate,gtfield=PurchaseDate" example:"2020-01-02T00:00:00Z"`
}
