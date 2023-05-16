package dto

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

type TokenDetails struct {
	Token     string
	UUID      string
	User      *TokenPayload
	ExpiresIn int64
}
