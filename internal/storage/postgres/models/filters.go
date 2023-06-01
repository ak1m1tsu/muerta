package models

type PageFilter struct {
	Limit  int
	Offset int
}

type ProductCategoryFilter struct {
	PageFilter
	Name string
}

type MeasureFilter struct {
	PageFilter
	Name string
}

type RecipeFilter struct {
	PageFilter
	Name string
}

type StepFilter struct {
	PageFilter
	Name string
}

type RoleFilter struct {
	PageFilter
	Name string
}

type ShelfLifeStatusFilter struct {
	PageFilter
	Name string
}

type ShelfLifeFilter struct {
	PageFilter
}

type StorageTypeFilter struct {
	PageFilter
	Name string
}

type StorageFilter struct {
	PageFilter
	Name string
}

type TipFilter struct {
	PageFilter
	Description string
}

type UserFilter struct {
	PageFilter
	Name string
}

type SettingFilter struct {
	PageFilter
	Name string
}
type ProductFilter struct {
	PageFilter
	Name string
}
