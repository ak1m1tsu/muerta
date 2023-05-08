package category

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type CategoryRepositorer interface {
	FindByID(ctx context.Context, id int) (models.ProductCategory, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.ProductCategory, error)
	Create(ctx context.Context, role models.ProductCategory) error
	Update(ctx context.Context, role models.ProductCategory) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
}

type categoryRepository struct {
	client repositories.PostgresClient
}

// Create implements CategoryRepositorer
func (r *categoryRepository) Create(ctx context.Context, role models.ProductCategory) error {
	var query = `
		INSERT INTO categories (name)
		VALUES ($1)
	`
	if _, err := r.client.Exec(ctx, query, role.Name); err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// Delete implements CategoryRepositorer
func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	var query = `
		UPDATE categories
		SET deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
	`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// FindByID implements CategoryRepositorer
func (r *categoryRepository) FindByID(ctx context.Context, id int) (models.ProductCategory, error) {
	var (
		query = `
			SELECT id, name, created_at
			FROM categories
			WHERE id = $1
			LIMIT 1
		`
		category models.ProductCategory
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name, &category.CreatedAt); err != nil {
		return models.ProductCategory{}, fmt.Errorf("failed to find category: %w", err)
	}
	return category, nil
}

// FindMany implements CategoryRepositorer
func (r *categoryRepository) FindMany(ctx context.Context, limit int, offset int, name string) ([]models.ProductCategory, error) {
	var (
		query = `
			SELECT id, name, created_at
			FROM categories
			WHERE name ILIKE $3
			ORDER BY created_at DESC
			LIMIT $1
			OFFSET $2
		`
		categories = make([]models.ProductCategory, 0, limit)
	)
	rows, err := r.client.Query(ctx, query, limit, offset, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var category models.ProductCategory
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// Restore implements CategoryRepositorer
func (r *categoryRepository) Restore(ctx context.Context, id int) error {
	var query = `
		UPDATE categories
		SET deleted_at = NULL,
			updated_at = NOW()
		WHERE id = $1
	`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to restore category: %w", err)
	}
	return nil
}

// Update implements CategoryRepositorer
func (r *categoryRepository) Update(ctx context.Context, role models.ProductCategory) error {
	var query = `
		UPDATE categories
		SET name = $1,
			updated_at = NOW()
		WHERE id = $2
	`
	if _, err := r.client.Exec(ctx, query, role.Name, role.ID); err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func New(client repositories.PostgresClient) CategoryRepositorer {
	return &categoryRepository{
		client: client,
	}
}
