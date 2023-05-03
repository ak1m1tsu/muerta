package dto

type LoginUserPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignUpUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserPayload struct {
	SignUpUserPayload
	Hash string
}

type TokenPayload struct {
	Name  string
	Roles []string
}

type FindStepDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Place int    `json:"place"`
}

type Paging struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
