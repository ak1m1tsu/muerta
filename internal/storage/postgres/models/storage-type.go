package models

type StorageType struct {
	ID   int    `db:"id,id_type"`
	Name string `db:"name"`
}
