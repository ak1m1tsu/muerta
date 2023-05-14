package dto

type CreateStepDTO struct {
	Name string `json:"name" validate:"required,gt=3,notblank"`
}

type UpdateStepDTO struct {
	Name string `json:"name" validate:"required,gt=3,notblank"`
}

type StepFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,notblank"`
}

func (f *StepFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}
func (f *StepFilterDTO) GetLimit() int {
	return f.Limit
}
func (f *StepFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}
func (f *StepFilterDTO) GetOffset() int {
	return f.Offset
}
