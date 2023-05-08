package storage

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type StorageRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Storage, error)
	FindMany(ctx context.Context, limit, offset int) ([]models.Storage, error)
	Create(ctx context.Context, storage *models.Storage) error
	Update(ctx context.Context, storage *models.Storage) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
}

type storageRepository struct {
	client repositories.PostgresClient
}

func New(client repositories.PostgresClient) StorageRepositorer {
	return &storageRepository{
		client: client,
	}
}

func (r *storageRepository) FindByID(ctx context.Context, id int) (models.Storage, error) {
	var (
		query = `
			SELECT 
				s.id, s.name, 
				s.temperature, s.humidity,
				s.created_at, st.id, st.name
			FROM storages s
			JOIN storages_types st ON st.id = s.id_type
			WHERE s.id = $1
			LIMIT 1
		`
		storage models.Storage
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&storage.ID, &storage.Name, &storage.Temperature, &storage.Humidity, &storage.CreatedAt, &storage.Type.ID, &storage.Type.Name); err != nil {
		return models.Storage{}, fmt.Errorf("find storage by id: %w", err)
	}
	return storage, nil
}

func (r *storageRepository) FindMany(ctx context.Context, limit, offset int) ([]models.Storage, error) {
	var (
		query = `
			SELECT 
				s.id, s.name, 
				s.temperature, s.humidity,
				s.created_at, st.name
			FROM storages s
			JOIN storages_types st ON st.id = s.id_type
			ORDER BY s.created_at DESC
			LIMIT $1 OFFSET $2
		`
		storages []models.Storage
	)
	rows, err := r.client.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("find many storages: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var storage models.Storage
		if err := rows.Scan(&storage.ID, &storage.Name, &storage.Temperature, &storage.Humidity, &storage.CreatedAt, &storage.Type.Name); err != nil {
			return nil, fmt.Errorf("scan storage: %w", err)
		}
		storages = append(storages, storage)
	}
	return storages, nil
}

func (r *storageRepository) Create(ctx context.Context, storage *models.Storage) error {
	var (
		query = `
			INSERT INTO storages 
				(name, temperature, humidity, id_type)
			VALUES
				($1, $2, $3, $4)
			RETURNING id
		`
	)
	if _, err := r.client.Exec(ctx, query, storage.Name, storage.Temperature, storage.Humidity, storage.Type.ID); err != nil {
		return fmt.Errorf("create storage: %w", err)
	}
	return nil
}

func (r *storageRepository) Update(ctx context.Context, storage *models.Storage) error {
	var (
		query = `
			UPDATE storages
			SET name = $1,
				temperature = $2,
				humidity = $3,
				id_type = $4,
				updated_at = NOW()
			WHERE id = $5
		`
	)
	if _, err := r.client.Exec(ctx, query, storage.Name, storage.Temperature, storage.Humidity, storage.Type.ID, storage.ID); err != nil {
		return fmt.Errorf("update storage: %w", err)
	}
	return nil
}

func (r *storageRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE storages
			SET deleted_at = NOW(),
				updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("delete storage: %w", err)
	}
	return nil
}

func (r *storageRepository) Restore(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE storages
			SET deleted_at = NULL,
				updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("restore storage: %w", err)
	}
	return nil
}
