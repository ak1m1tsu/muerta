package models

import "time"

type ShelfLife struct {
	ID           int `db:"id"`
	Product      Product
	Storage      Storage
	Measure      Measure
	User         User
	Quantity     int        `db:"quantity"`
	PurchaseDate *time.Time `db:"purchase_date"`
	EndDate      *time.Time `db:"end_date"`
	CreatedAt    *time.Time `db:"created_at"`
}

type ShelfLifeStatus struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
