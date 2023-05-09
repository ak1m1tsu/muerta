package tip

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type TipRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Tip, error)
	FindMany(ctx context.Context, limit, offset int, description string) ([]models.Tip, error)
	Create(ctx context.Context, tip models.Tip) error
	Update(ctx context.Context, tip models.Tip) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
}

type tipRepository struct {
	client repositories.PostgresClient
}

// Create implements TipRepositorer
func (r *tipRepository) Create(ctx context.Context, tip models.Tip) error {
	var (
		query = `
			INSERT INTO tips (description)
			VALUES ($1)
		`
	)
	if _, err := r.client.Exec(ctx, query, tip.Description); err != nil {
		return fmt.Errorf("failed to create tip: %w", err)
	}
	return nil
}

// Delete implements TipRepositorer
func (r *tipRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE tips
			SET deleted_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete tip: %w", err)
	}
	return nil
}

// FindByID implements TipRepositorer
func (r *tipRepository) FindByID(ctx context.Context, id int) (models.Tip, error) {
	var (
		query = `
			SELECT id, description
			FROM tips
			WHERE id = $1
			LIMIT 1	
		`
		tip models.Tip
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&tip.ID, &tip.Description); err != nil {
		return models.Tip{}, fmt.Errorf("failed to find tip: %w", err)
	}
	return tip, nil
}

// FindMany implements TipRepositorer
func (r *tipRepository) FindMany(ctx context.Context, limit int, offset int, description string) ([]models.Tip, error) {
	var (
		query = `
			SELECT id, description
			FROM tips
			WHERE description ILIKE $1
			LIMIT $2
			OFFSET $3
		`
		tips = make([]models.Tip, 0, limit)
	)
	rows, err := r.client.Query(ctx, query, "%"+description+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find tips: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var tip models.Tip
		if err := rows.Scan(&tip.ID, &tip.Description); err != nil {
			return nil, fmt.Errorf("failed to scan tip: %w", err)
		}
		tips = append(tips, tip)
	}
	return tips, nil
}

// Restore implements TipRepositorer
func (r *tipRepository) Restore(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE tips
			SET deleted_at = NULL
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to restore tip: %w", err)
	}
	return nil
}

// Update implements TipRepositorer
func (r *tipRepository) Update(ctx context.Context, tip models.Tip) error {
	var (
		query = `
			UPDATE tips
			SET description = $1
			WHERE id = $2
		`
	)
	if _, err := r.client.Exec(ctx, query, tip.Description, tip.ID); err != nil {
		return fmt.Errorf("failed to update tip: %w", err)
	}
	return nil
}

func New(client repositories.PostgresClient) TipRepositorer {
	return &tipRepository{
		client: client,
	}
}
