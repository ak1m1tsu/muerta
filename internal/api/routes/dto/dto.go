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
