package dto

type CreateRole struct {
	Name string `json:"name" validate:"required,gte=3,lte=20,alpha" example:"admin"`
}

type UpdateRole struct {
	Name string `json:"name" validate:"required,gte=3,lte=20,alpha" example:"admin"`
}

type FindRole struct {
	ID   int    `json:"id"   example:"1"`
	Name string `json:"name" example:"admin"`
}
