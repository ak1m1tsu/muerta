package params

type Paging struct {
	Limit  int `query:"limit"  example:"10" validate:"required,oneof=5 10 15 20 25 30"`
	Offset int `query:"offset" example:"0"  validate:"omitempty,gte=0"`
}

type ProductCategoryFilter struct {
	Paging
	Name string `query:"name" example:"овощь" validate:"omitempty,gte=1,notblank"`
}

type ProductFilter struct {
	Paging
	Name string `query:"name" validate:"omitempty,gte=1,notblank" example:"помидор"`
}

type MeasureFilter struct {
	Paging
	Name string `query:"name" example:"кг" validate:"omitempty,gte=1,notblank"`
}

type RecipeFilter struct {
	Paging
	Name string `query:"name" example:"салат" validate:"omitempty,gte=1,notblank"`
}

type StepFilter struct {
	Paging
	Name string `query:"name" example:"налить воду в емкость" validate:"omitempty,gte=1,notblank"`
}

type RoleFilter struct {
	Paging
	Name string `query:"name" example:"user" validate:"omitempty,gte=1,alpha"`
}

type ShelfLifeStatusFilter struct {
	Paging
	Name string `query:"name" example:"просрочен" validate:"omitempty,gte=1,notblank"`
}

type ShelfLifeFilter struct {
	Paging
}

type StorageTypeFilter struct {
	Paging
	Name string `query:"name" example:"хрупкое" validate:"omitempty,gte=1,notblank"`
}

type StorageFilter struct {
	Paging
	Name string `query:"name" example:"банка" validate:"omitempty,gte=1,notblank"`
}

type TipFilter struct {
	Paging
	Description string `query:"description" example:"хранить при низкой температуре" validate:"omitempty,gte=1,notblank"`
}

type UserFilter struct {
	Paging
	Name string `query:"name" example:"hunter" validate:"omitempty,gte=1,alpha"`
}

type SettingFilter struct {
	Paging
	Name string `query:"name" example:"получать рассылку" validate:"omitempty,gte=1,notblank"`
}
