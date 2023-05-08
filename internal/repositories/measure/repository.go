package measure

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type MeasureRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Measure, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.Measure, error)
	Create(ctx context.Context, role models.Measure) error
	Update(ctx context.Context, role models.Measure) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
}

type measureRepository struct {
	client repositories.PostgresClient
}

// Create implements MeasureRepositorer
func (r *measureRepository) Create(ctx context.Context, role models.Measure) error {
	var (
		query = `
			INSERT INTO roles (name)
			VALUES ($1)
		`
	)
	if _, err := r.client.Exec(ctx, query, role.Name); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// Delete implements MeasureRepositorer
func (r *measureRepository) Delete(ctx context.Context, id int) error {
	panic("unimplemented")
}

// FindByID implements MeasureRepositorer
func (r *measureRepository) FindByID(ctx context.Context, id int) (models.Measure, error) {
	panic("unimplemented")
}

// FindMany implements MeasureRepositorer
func (r *measureRepository) FindMany(ctx context.Context, limit int, offset int, name string) ([]models.Measure, error) {
	panic("unimplemented")
}

// Restore implements MeasureRepositorer
func (r *measureRepository) Restore(ctx context.Context, id int) error {
	panic("unimplemented")
}

// Update implements MeasureRepositorer
func (r *measureRepository) Update(ctx context.Context, role models.Measure) error {
	panic("unimplemented")
}

func New(client repositories.PostgresClient) MeasureRepositorer {
	return &measureRepository{
		client: client,
	}
}
