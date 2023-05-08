package dto

type RoleFilterDTO struct {
	*Paging
	Name string `query:"name"`
}
type CreateRoleDTO struct {
	Name string `json:"name" validate:"required"`
}
type UpdateRoleDTO struct {
	Name string `json:"name" validate:"required"`
}
type FindRoleDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
