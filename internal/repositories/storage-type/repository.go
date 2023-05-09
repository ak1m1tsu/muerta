package storagetype

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type StorageTypeRepositorer interface {
	FindByID(ctx context.Context, id int) (models.StorageType, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.StorageType, error)
	Create(ctx context.Context, storageType models.StorageType) error
	Update(ctx context.Context, storageType models.StorageType) error
	Delete(ctx context.Context, id int) error
}

type storageTypeRepository struct {
	client repositories.PostgresClient
}

// Create implements StorageTypeRepositorer
func (r *storageTypeRepository) Create(ctx context.Context, storageType models.StorageType) error {
	var (
		query = `
			INSERT INTO storages_types (name)
			VALUES ($1)
		`
	)
	if _, err := r.client.Exec(ctx, query, storageType.Name); err != nil {
		return fmt.Errorf("failed to create storageType: %w", err)
	}
	return nil
}

// Delete implements StorageTypeRepositorer
func (r *storageTypeRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			DELETE FROM storages_types
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete storageType: %w", err)
	}
	return nil
}

// FindByID implements StorageTypeRepositorer
func (r *storageTypeRepository) FindByID(ctx context.Context, id int) (models.StorageType, error) {
	var (
		query = `
			SELECT id, name
			FROM storages_types
			WHERE id = $1
			LIMIT 1	
		`
		storageType models.StorageType
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&storageType.ID, &storageType.Name); err != nil {
		return models.StorageType{}, fmt.Errorf("failed to find storageType: %w", err)
	}
	return storageType, nil
}

// FindMany implements StorageTypeRepositorer
func (r *storageTypeRepository) FindMany(ctx context.Context, limit int, offset int, name string) ([]models.StorageType, error) {
	var (
		query = `
			SELECT id, name
			FROM storages_types
			WHERE name ILIKE $1
			LIMIT $2
			OFFSET $3
		`
		storageTypes = make([]models.StorageType, 0, limit)
	)
	rows, err := r.client.Query(ctx, query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find storageTypes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var storageType models.StorageType
		if err := rows.Scan(&storageType.ID, &storageType.Name); err != nil {
			return nil, fmt.Errorf("failed to scan storageType: %w", err)
		}
		storageTypes = append(storageTypes, storageType)
	}
	return storageTypes, nil
}

// Update implements StorageTypeRepositorer
func (r *storageTypeRepository) Update(ctx context.Context, storageType models.StorageType) error {
	var (
		query = `
			UPDATE storages_types
			SET name = $1
			WHERE id = $2
		`
	)
	if _, err := r.client.Exec(ctx, query, storageType.Name, storageType.ID); err != nil {
		return fmt.Errorf("failed to update storageType: %w", err)
	}
	return nil
}

func New(client repositories.PostgresClient) StorageTypeRepositorer {
	return &storageTypeRepository{
		client: client,
	}
}
