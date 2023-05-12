package dto

type TokenDetails struct {
	Token     string
	UUID      string
	User      *TokenPayload
	ExpiresIn int64
}
