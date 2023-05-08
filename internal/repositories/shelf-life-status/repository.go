package shelflifestatus

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type ShelfLifeStatusRepositorer interface {
	FindByID(ctx context.Context, id int) (models.ShelfLifeStatus, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.ShelfLifeStatus, error)
	Create(ctx context.Context, shelfLifeStatus models.ShelfLifeStatus) error
	Update(ctx context.Context, shelfLifeStatus models.ShelfLifeStatus) error
	Delete(ctx context.Context, id int) error
}

type shelfLifeStatusRepository struct {
	client repositories.PostgresClient
}

// Create implements ShelfLifeStatusRepositorer
func (r *shelfLifeStatusRepository) Create(ctx context.Context, shelfLifeStatus models.ShelfLifeStatus) error {
	var (
		query = `
			INSERT INTO statuses (name)
			VALUES ($1)
		`
	)
	if _, err := r.client.Exec(ctx, query, shelfLifeStatus.Name); err != nil {
		return fmt.Errorf("failed to create shelfLifeStatus: %w", err)
	}
	return nil
}

// Delete implements ShelfLifeStatusRepositorer
func (r *shelfLifeStatusRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			DELETE FROM statuses
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete shelfLifeStatus: %w", err)
	}
	return nil
}

// FindByID implements ShelfLifeStatusRepositorer
func (r *shelfLifeStatusRepository) FindByID(ctx context.Context, id int) (models.ShelfLifeStatus, error) {
	var (
		query = `
			SELECT id, name
			FROM statuses
			WHERE id = $1
			LIMIT 1	
		`
		shelfLifeStatus models.ShelfLifeStatus
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&shelfLifeStatus.ID, &shelfLifeStatus.Name); err != nil {
		return models.ShelfLifeStatus{}, fmt.Errorf("failed to find shelfLifeStatus: %w", err)
	}
	return shelfLifeStatus, nil
}

// FindMany implements ShelfLifeStatusRepositorer
func (r *shelfLifeStatusRepository) FindMany(ctx context.Context, limit int, offset int, name string) ([]models.ShelfLifeStatus, error) {
	var (
		query = `
			SELECT id, name
			FROM statuses
			WHERE name ILIKE $1
			LIMIT $2
			OFFSET $3
		`
		statuses = make([]models.ShelfLifeStatus, 0, limit)
	)
	rows, err := r.client.Query(ctx, query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find statuses: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var shelfLifeStatus models.ShelfLifeStatus
		if err := rows.Scan(&shelfLifeStatus.ID, &shelfLifeStatus.Name); err != nil {
			return nil, fmt.Errorf("failed to scan shelfLifeStatus: %w", err)
		}
		statuses = append(statuses, shelfLifeStatus)
	}
	return statuses, nil
}

// Update implements ShelfLifeStatusRepositorer
func (r *shelfLifeStatusRepository) Update(ctx context.Context, shelfLifeStatus models.ShelfLifeStatus) error {
	var (
		query = `
			UPDATE statuses
			SET name = $1
			WHERE id = $2
		`
	)
	if _, err := r.client.Exec(ctx, query, shelfLifeStatus.Name, shelfLifeStatus.ID); err != nil {
		return fmt.Errorf("failed to update shelfLifeStatus: %w", err)
	}
	return nil
}

func New(client repositories.PostgresClient) ShelfLifeStatusRepositorer {
	return &shelfLifeStatusRepository{
		client: client,
	}
}
