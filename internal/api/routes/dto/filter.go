package dto

type Paging struct {
	Limit  int `query:"limit" example:"10" validate:"required,oneof=5 10 15 20 25 30"`
	Offset int `query:"offset" example:"0" validate:"omitempty,gte=0"`
}

type ProductCategoryFilterDTO struct {
	Paging
	Name string `query:"name" example:"овощь" validate:"omitempty,gte=1,notblank"`
}

type ProductFilterDTO struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,notblank"`
}

type MeasureFilterDTO struct {
	Paging
	Name string `query:"name" example:"кг" validate:"omitempty,gte=1,notblank"`
}

type RecipeFilterDTO struct {
	Paging
	Name string `query:"name" example:"салат" validate:"omitempty,gte=1,notblank"`
}

type StepFilterDTO struct {
	Paging
	Name string `query:"name" example:"налить воду в емкость" validate:"omitempty,gte=1,notblank"`
}

type RoleFilterDTO struct {
	Paging
	Name string `query:"name" example:"user" validate:"omitempty,gte=1,alpha"`
}

type ShelfLifeStatusFilterDTO struct {
	Paging
	Name string `query:"name" example:"просрочен" validate:"omitempty,gte=1,notblank"`
}

type ShelfLifeFilterDTO struct {
	Paging
}

type StorageTypeFilterDTO struct {
	Paging
	Name string `query:"name" example:"хрупкое" validate:"omitempty,gte=1,notblank"`
}

type StorageFilterDTO struct {
	Paging
	Name string `query:"name" example:"банка" validate:"omitempty,gte=1,notblank"`
}

type TipFilterDTO struct {
	Paging
	Description string `query:"description" example:"хранить при низкой температуре" validate:"omitempty,gte=1,notblank"`
}

type UserFilterDTO struct {
	Paging
	Name string `query:"name" example:"hunter" validate:"omitempty,gte=1,alpha"`
}

type SettingFilterDTO struct {
	Paging
	Name string `query:"name" example:"получать рассылку" validate:"omitempty,gte=1,notblank"`
}
