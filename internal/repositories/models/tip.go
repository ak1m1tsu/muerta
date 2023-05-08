package models

type Tip struct {
	ID          int    `db:"id,id_tip"`
	Description string `db:"description"`
}
