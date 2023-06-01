package role

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

type RoleRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Role, error)
	FindByName(ctx context.Context, name string) (models.Role, error)
	FindMany(ctx context.Context, filter models.RoleFilter) ([]models.Role, error)
	Create(ctx context.Context, role models.Role) error
	Update(ctx context.Context, role models.Role) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	Count(ctx context.Context, filter models.RoleFilter) (int, error)
}

type roleRepository struct {
	client postgres.Client
}

// FindByName implements RoleRepositorer
func (r *roleRepository) FindByName(ctx context.Context, name string) (models.Role, error) {
	var (
		query = `
			SELECT id, name FROM roles
			WHERE name = $1
			LIMIT 1
		`
		role models.Role
	)
	if err := r.client.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name); err != nil {
		return role, fmt.Errorf("failed to find role by name: %w", err)
	}
	return role, nil
}

func New(client postgres.Client) RoleRepositorer {
	return &roleRepository{
		client: client,
	}
}

func (r *roleRepository) Count(ctx context.Context, filter models.RoleFilter) (int, error) {
	var (
		query = `
			SELECT COUNT(*) 
			FROM roles 
			WHERE deleted_at IS NULL AND
				name ILIKE $1
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query, "%"+filter.Name+"%").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count roles: %w", err)
	}
	return count, nil
}

// Create implements RoleRepositorer
func (r *roleRepository) Create(ctx context.Context, role models.Role) error {
	query := `
			INSERT INTO roles (name)
			VALUES ($1)
		`
	if _, err := r.client.Exec(ctx, query, role.Name); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// Delete implements RoleRepositorer
func (r *roleRepository) Delete(ctx context.Context, id int) error {
	query := `
			UPDATE roles
			SET deleted_at = NOW()
			WHERE id = $1
		`
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
func (r *roleRepository) FindMany(
	ctx context.Context,
	filter models.RoleFilter,
) ([]models.Role, error) {
	var (
		query = `
			SELECT id, name
			FROM roles
			WHERE name ILIKE $1 AND
				deleted_at IS NULL
			LIMIT $2
			OFFSET $3
		`
		roles = make([]models.Role, 0, filter.Limit)
	)
	rows, err := r.client.Query(ctx, query, "%"+filter.Name+"%", filter.Limit, filter.Offset)
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
	query := `
			UPDATE roles
			SET deleted_at = NULL
			WHERE id = $1
		`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to restore role: %w", err)
	}
	return nil
}

// Update implements RoleRepositorer
func (r *roleRepository) Update(ctx context.Context, role models.Role) error {
	query := `
			UPDATE roles
			SET name = $1
			WHERE id = $2
		`
	if _, err := r.client.Exec(ctx, query, role.Name, role.ID); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}
