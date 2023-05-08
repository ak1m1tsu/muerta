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
		*UserFilterDTO

	SetLimit(int)
	GetLimit() int
	SetOffset(int)
	GetOffset() int
}

type LoginUserPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignUpUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type TokenPayload struct {
	Name  string
	Roles []string
}

type Paging struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
