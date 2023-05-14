package dto

type MeasureFilterDTO struct {
	Paging
	Name string `query:"name" example:"кг" validate:"omitempty,gte=1,notblank"`
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
	Name string `json:"name" validate:"required,gte=1,notblank" example:"кг"`
}

type UpdateMeasureDTO struct {
	Name string `json:"name" validate:"required,gte=1,notblank" example:"л"`
}

type FindMeasureDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
