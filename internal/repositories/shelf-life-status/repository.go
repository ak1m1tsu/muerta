package shelflifestatus

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type ShelfLifeStatusRepositorer interface {
	FindByID(ctx context.Context, id int) (models.ShelfLifeStatus, error)
	FindMany(
		ctx context.Context,
		filter models.ShelfLifeStatusFilter,
	) ([]models.ShelfLifeStatus, error)
	Create(ctx context.Context, shelfLifeStatus models.ShelfLifeStatus) error
	Update(ctx context.Context, shelfLifeStatus models.ShelfLifeStatus) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context, filter models.ShelfLifeStatusFilter) (int, error)
}

type shelfLifeStatusRepository struct {
	client repositories.PostgresClient
}

func (r *shelfLifeStatusRepository) Count(
	ctx context.Context,
	filter models.ShelfLifeStatusFilter,
) (int, error) {
	var (
		query = `
			SELECT COUNT(*) 
			FROM statuses
			WHERE name ILIKE $1
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query, "%"+filter.Name+"%").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count statuses: %w", err)
	}
	return count, nil
}

// Create implements ShelfLifeStatusRepositorer
func (r *shelfLifeStatusRepository) Create(
	ctx context.Context,
	shelfLifeStatus models.ShelfLifeStatus,
) error {
	query := `
			INSERT INTO statuses (name)
			VALUES ($1)
		`
	if _, err := r.client.Exec(ctx, query, shelfLifeStatus.Name); err != nil {
		return fmt.Errorf("failed to create shelfLifeStatus: %w", err)
	}
	return nil
}

// Delete implements ShelfLifeStatusRepositorer
func (r *shelfLifeStatusRepository) Delete(ctx context.Context, id int) error {
	query := `
			DELETE FROM statuses
			WHERE id = $1
		`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete shelfLifeStatus: %w", err)
	}
	return nil
}

// FindByID implements ShelfLifeStatusRepositorer
func (r *shelfLifeStatusRepository) FindByID(
	ctx context.Context,
	id int,
) (models.ShelfLifeStatus, error) {
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
func (r *shelfLifeStatusRepository) FindMany(
	ctx context.Context,
	filter models.ShelfLifeStatusFilter,
) ([]models.ShelfLifeStatus, error) {
	var (
		query = `
			SELECT id, name
			FROM statuses
			WHERE name ILIKE $1
			LIMIT $2
			OFFSET $3
		`
		statuses = make([]models.ShelfLifeStatus, 0, filter.Limit)
	)
	rows, err := r.client.Query(ctx, query, "%"+filter.Name+"%", filter.Limit, filter.Offset)
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
func (r *shelfLifeStatusRepository) Update(
	ctx context.Context,
	shelfLifeStatus models.ShelfLifeStatus,
) error {
	query := `
			UPDATE statuses
			SET name = $1
			WHERE id = $2
		`
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
