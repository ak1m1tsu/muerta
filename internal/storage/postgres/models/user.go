package models

import "time"

type User struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Salt      string `db:"salt"`
	Settings  []Setting
	Roles     []Role
	Password  Password
	CreatedAt time.Time `db:"created_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

type Setting struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Value    string `db:"value"`
	Category SettingCategory
}

type SettingCategory struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Password struct {
	Hash string `db:"passhash"`
}
