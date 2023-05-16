package dto

type CreateSettingDTO struct {
	Name       string `json:"name,omitempty" validate:"required,gte=3,alphanumunicode"`
	Value      string `json:"value,omitempty" validate:"required,gt=0,alphanumunicode"`
	CategoryID int    `json:"id_category,omitempty" validate:"required,gt=0"`
}

type UpdateSettingDTO struct {
	Name       string `json:"name" validate:"omitempty,gte=3,notblank"`
	CategoryID int    `json:"id_category" validate:"omitempty,gt=0"`
}

type FindSettingDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    string `json:"value,omitempty"`
	Category string `json:"category"`
}
