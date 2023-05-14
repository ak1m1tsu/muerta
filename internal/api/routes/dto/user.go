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
	Value string `json:"value" validate:"required,gt=0,alphanumunicode"`
}

type UserRoleDTO struct {
	ID int `json:"id"`
}

type CreateSettingDTO struct {
	Name       string `json:"name,omitempty" validate:"required,gte=3,alphanumunicode"`
	Value      string `json:"value,omitempty" validate:"required,gt=0,alphanumunicode"`
	CategoryID int    `json:"id_category,omitempty" validate:"required,gt=0"`
}

type UpdateSettingDTO struct {
	Name       string `json:"name" validate:"omitempty,gte=3,alphanumunicode"`
	CategoryID int    `json:"id_category" validate:"omitempty,gt=0"`
}

type UpdateUserSettingDTO struct {
	Value string `json:"value" validate:"required,gt=0,alphanumunicode"`
}

type FindSettingDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value,omitempty"`
	Category string `json:"category"`
}

type UserFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,alphaunicode"`
}

func (f *UserFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *UserFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *UserFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *UserFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type SettingFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,alphaunicode"`
}

func (f *SettingFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *SettingFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *SettingFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *SettingFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type UserStorageDTO struct {
	StorageID int `json:"id_storage" validate:"gt=0"`
}

type UserShelfLifeDTO struct {
	ShelfLifeID  int        `json:"id_shelf_life" validate:"gt=0"`
	ProductID    int        `json:"id_product"`
	StorageID    int        `json:"id_storage"`
	MeasureID    int        `json:"id_measure"`
	Quantity     int        `json:"quantity"`
	PurchaseDate *time.Time `json:"purchase_date"`
	EndDate      *time.Time `json:"end_date"`
}
