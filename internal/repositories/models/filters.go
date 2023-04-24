package models

import "time"

type BaseFilter struct {
	Limit  int
	Offset int
}

type UserFilter struct {
	BaseFilter
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
