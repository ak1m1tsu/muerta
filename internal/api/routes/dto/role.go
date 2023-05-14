package dto

type RoleFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,alpha"`
}
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

func (f *RoleFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *RoleFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *RoleFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *RoleFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}
