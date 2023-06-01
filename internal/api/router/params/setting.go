package params

type CreateSetting struct {
	Name       string `json:"name,omitempty"        validate:"required,gte=3,alphanumunicode" example:"Уведомления на почту"`
	Value      string `json:"value,omitempty"       validate:"required,gt=0,alphanumunicode"  example:"XXXXXXXXXXXX"`
	CategoryID int    `json:"id_category,omitempty" validate:"required,gt=0"                  example:"1"`
}

type UpdateSetting struct {
	Name       string `json:"name"        validate:"omitempty,gte=3,notblank" example:"Уведомления на почту"`
	CategoryID int    `json:"id_category" validate:"omitempty,gt=0"           example:"1"`
}

type FindSetting struct {
	ID       int    `json:"id"              example:"1"`
	Name     string `json:"name"            example:"Уведомления на почту"`
	Value    string `json:"value,omitempty" example:"XXXXXXXXXXXX"`
	Category string `json:"category"        example:"Уведомления"`
}
