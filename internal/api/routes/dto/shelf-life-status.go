package dto

type ShelfLifeStatusFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,alphaunicode"`
}

func (f *ShelfLifeStatusFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *ShelfLifeStatusFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *ShelfLifeStatusFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *ShelfLifeStatusFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type CreateShelfLifeStatusDTO struct {
	Name string `json:"name" validate:"required,gte=3,alphaunicode"`
}
type UpdateShelfLifeStatusDTO struct {
	Name string `json:"name" validate:"required,gte=3,alphaunicode"`
}
type FindShelfLifeStatusDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
