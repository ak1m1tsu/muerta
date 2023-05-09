package dto

type MeasureFilterDTO struct {
	Paging
	Name string `query:"name"`
}

func (f *MeasureFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *MeasureFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *MeasureFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *MeasureFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type CreateMeasureDTO struct {
	Name string `json:"name" validate:"required"`
}

type UpdateMeasureDTO struct {
	Name string `json:"name" validate:"required"`
}

type FindMeasureDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
