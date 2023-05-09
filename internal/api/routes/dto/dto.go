package dto

type Filter interface {
	*CategoryFilterDTO |
		*MeasureFilterDTO |
		*ProductFilterDTO |
		*RecipeFilterDTO |
		*RoleFilterDTO |
		*StorageFilterDTO |
		*SettingFilterDTO |
		*TipFilterDTO |
		*UserFilterDTO |
		*StorageTypeFilterDTO |
		*ShelfLifeFilterDTO |
		*ShelfLifeStatusFilterDTO

	SetLimit(int)
	GetLimit() int
	SetOffset(int)
	GetOffset() int
}

type LoginDTO struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignUpDTO struct {
	Name            string `json:"name" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8"`
}

type TokenPayload struct {
	ID    int
	Name  string
	Roles []interface{}
}

type Paging struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
