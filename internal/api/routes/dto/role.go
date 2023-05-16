package dto

type CreateRoleDTO struct {
	Name string `json:"name" validate:"required,gte=3,lte=20,alpha"`
}

type UpdateRoleDTO struct {
	Name string `json:"name" validate:"required,gte=3,lte=20,alpha"`
}

type FindRoleDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
