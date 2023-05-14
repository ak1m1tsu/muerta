package dto

type ProductFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,alphaunicode"`
}

func (f *ProductFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *ProductFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *ProductFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *ProductFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type CreateProductDTO struct {
	Name string `json:"name" validate:"required,gte=2,alphaunicode"`
}

type UpdateProductDTO struct {
	Name string `json:"name" validate:"required,gte=2,alphaunicode"`
}

type FindProductDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
