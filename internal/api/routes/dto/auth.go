package dto

type Login struct {
	Name     string `json:"name"     validate:"required,gte=3,alpha"    example:"theBestUserEver"`
	Password string `json:"password" validate:"required,gte=8,alphanum" example:"th3B3stUs3rEver"`
}

type SignUp struct {
	Name            string `json:"name"             validate:"required,gte=3,alpha"                     example:"theBestUserEver"`
	Password        string `json:"password"         validate:"required,gte=8,alphanum"                  example:"th3B3stUs3rEver"`
	PasswordConfirm string `json:"password_confirm" validate:"required,gte=8,alphanum,eqfield=Password" example:"th3B3stUs3rEver"`
}

type TokenPayload struct {
	UUID     string
	UserID   int
	Username string
	Roles    []string
}

type TokenDetails struct {
	Token     string
	UUID      string
	User      *TokenPayload
	ExpiresIn int64
}
