package models

import "time"

type Vault struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Temperature float32    `db:"temperature"`
	Humidity    float32    `db:"humidity"`
	CreatedAt   *time.Time `db:"created_at"`
	Type        StorageType
}
