package models

import (
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
)

func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		user.ID,
		user.Name,
		*user.CreatedAt,
	}
}

func FilterUserPayload(payload dto.CreateUserPayload) User {
	return User{
		Name: payload.Name,
	}
}
