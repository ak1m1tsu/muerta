package measure

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type MeasureRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Measure, error)
	FindMany(ctx context.Context, filter models.MeasureFilter) ([]models.Measure, error)
	Create(ctx context.Context, measure models.Measure) error
	Update(ctx context.Context, measure models.Measure) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context, filter models.MeasureFilter) (int, error)
}

type measureRepository struct {
	client repositories.PostgresClient
}

func (r *measureRepository) Count(ctx context.Context, filter models.MeasureFilter) (int, error) {
	var (
		query = `
			SELECT COUNT(*) FROM measures WHERE name ILIKE $1
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query, "%"+filter.Name+"%").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count measures: %w", err)
	}
	return count, nil
}

// Create implements MeasureRepositorer
func (r *measureRepository) Create(ctx context.Context, measure models.Measure) error {
	query := `
			INSERT INTO measures (name)
			VALUES ($1)
		`
	if _, err := r.client.Exec(ctx, query, measure.Name); err != nil {
		return fmt.Errorf("failed to create measure: %w", err)
	}
	return nil
}

// Delete implements MeasureRepositorer
func (r *measureRepository) Delete(ctx context.Context, id int) error {
	query := `
			DELETE FROM measures
			WHERE id = $1
		`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete measure: %w", err)
	}
	return nil
}

// FindByID implements MeasureRepositorer
func (r *measureRepository) FindByID(ctx context.Context, id int) (models.Measure, error) {
	var (
		query = `
			SELECT id, name
			FROM measures
			WHERE id = $1
			LIMIT 1	
		`
		measure models.Measure
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&measure.ID, &measure.Name); err != nil {
		return models.Measure{}, fmt.Errorf("failed to find measure: %w", err)
	}
	return measure, nil
}

// FindMany implements MeasureRepositorer
func (r *measureRepository) FindMany(
	ctx context.Context,
	filter models.MeasureFilter,
) ([]models.Measure, error) {
	var (
		query = `
			SELECT id, name
			FROM measures
			WHERE name ILIKE $1
			LIMIT $2
			OFFSET $3
		`
		measures = make([]models.Measure, 0, filter.Limit)
	)
	rows, err := r.client.Query(ctx, query, "%"+filter.Name+"%", filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find measures: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var measure models.Measure
		if err := rows.Scan(&measure.ID, &measure.Name); err != nil {
			return nil, fmt.Errorf("failed to scan measure: %w", err)
		}
		measures = append(measures, measure)
	}
	return measures, nil
}

// Update implements MeasureRepositorer
func (r *measureRepository) Update(ctx context.Context, measure models.Measure) error {
	query := `
			UPDATE measures
			SET name = $1
			WHERE id = $2
		`
	if _, err := r.client.Exec(ctx, query, measure.Name, measure.ID); err != nil {
		return fmt.Errorf("failed to update measure: %w", err)
	}
	return nil
}

func New(client repositories.PostgresClient) MeasureRepositorer {
	return &measureRepository{
		client: client,
	}
}
