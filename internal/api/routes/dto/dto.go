package dto

type Filter interface {
	*ProductCategoryFilterDTO |
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
		*ShelfLifeStatusFilterDTO |
		*StepFilterDTO

	SetLimit(int)
	GetLimit() int
	SetOffset(int)
	GetOffset() int
}

type LoginDTO struct {
	Name     string `json:"name" validate:"required,gte=3,alpha"`
	Password string `json:"password" validate:"required,gte=8,alphanum"`
}

type SignUpDTO struct {
	Name            string `json:"name" validate:"required,gte=3,alpha"`
	Password        string `json:"password" validate:"required,gte=8,alphanum"`
	PasswordConfirm string `json:"password_confirm" validate:"required,gte=8,alphanum,eqfield=Password"`
}

type TokenPayload struct {
	ID    int
	Name  string
	Roles []interface{}
}

type Paging struct {
	Limit  int `query:"limit" example:"10" validate:"omitempty,gte=0"`
	Offset int `query:"offset" example:"0" validate:"omitempty,gte=0"`
}
