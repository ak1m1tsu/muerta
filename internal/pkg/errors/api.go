package errors

var (
	ErrNotAdmin = New("user is not admin")
	ErrNotOwner = New("user is not owner")
)

var (
	ErrFailedToGetTokenPayload = New("failed to get token payload")
	ErrFailedToGetUserId       = New("failed to get user id")
)
