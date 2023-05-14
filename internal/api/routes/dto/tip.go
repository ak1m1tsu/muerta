package dto

type FindTipDTO struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type UpdateTipDTO struct {
	Description string `json:"description" validate:"required,gte=3,lte=200"`
}

type CreateTipDTO struct {
	Description string `json:"description" validate:"required,gte=3,lte=200"`
}

type TipFilterDTO struct {
	Paging
	Description string `query:"description" validate:"omitempty,gte=1,alphaunicode"`
}

func (f *TipFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *TipFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *TipFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *TipFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}
