package role

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type RoleRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Role, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.Role, error)
	Create(ctx context.Context, role models.Role) error
	Update(ctx context.Context, role models.Role) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
}

type roleRepository struct {
	client repositories.PostgresClient
}

func New(client repositories.PostgresClient) RoleRepositorer {
	return &roleRepository{
		client: client,
	}
}

// Create implements RoleRepositorer
func (r *roleRepository) Create(ctx context.Context, role models.Role) error {
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

// Delete implements RoleRepositorer
func (r *roleRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE roles
			SET deleted_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// FindByID implements RoleRepositorer
func (r *roleRepository) FindByID(ctx context.Context, id int) (models.Role, error) {
	var (
		query = `
			SELECT id, name
			FROM roles
			WHERE id = $1
			LIMIT 1	
		`
		role models.Role
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name); err != nil {
		return models.Role{}, fmt.Errorf("failed to find role: %w", err)
	}
	return role, nil
}

// FindMany implements RoleRepositorer
func (r *roleRepository) FindMany(ctx context.Context, limit int, offset int, name string) ([]models.Role, error) {
	var (
		query = `
			SELECT id, name
			FROM roles
			WHERE name ILIKE $1
			LIMIT $2
			OFFSET $3
		`
		roles = make([]models.Role, 0, limit)
	)
	rows, err := r.client.Query(ctx, query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find roles: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// Restore implements RoleRepositorer
func (r *roleRepository) Restore(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE roles
			SET deleted_at = NULL
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to restore role: %w", err)
	}
	return nil
}

// Update implements RoleRepositorer
func (r *roleRepository) Update(ctx context.Context, role models.Role) error {
	var (
		query = `
			UPDATE roles
			SET name = $1
			WHERE id = $2
		`
	)
	if _, err := r.client.Exec(ctx, query, role.Name, role.ID); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}