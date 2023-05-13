package dto

import "time"

type UserPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type FindUserDTO struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	CreatedAt time.Time        `json:"created_at,omitempty"`
	Settings  []FindSettingDTO `json:"settings,omitempty"`
}

type UpdateUserDTO struct {
	Name    string `json:"name"`
	Restore bool   `json:"restore"`
}

type CreateUserDTO struct {
	ID       int                `json:"_"`
	Name     string             `json:"name" validate:"required"`
	Password string             `json:"password" validate:"required,min=8"`
	Settings []CreateSettingDTO `json:"settings"`
}

type CreateSettingDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	Value      string `json:"value,omitempty"`
	CategoryID int    `json:"id_category,omitempty"`
}

type UpdateSettingDTO struct {
	Name       string `json:"name"`
	CategoryID int    `json:"id_category"`
}

type UpdateUserSettingDTO struct {
	Value string `json:"value" validate:"required"`
}

type FindSettingDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value,omitempty"`
	Category string `json:"category"`
}

type UserFilterDTO struct {
	Paging
	Name string `query:"name"`
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
	Name string `query:"name"`
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
