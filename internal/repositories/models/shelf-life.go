package models

import "time"

type ShelfLife struct {
	ID           int `db:"id"`
	Product      Product
	Storage      Storage
	Measure      Measure
	Quantity     int        `db:"quantity"`
	PurchaseDate *time.Time `db:"purchase_date"`
	EndDate      *time.Time `db:"end_date"`
	CreatedAt    *time.Time `db:"created_at"`
}
