package setting

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type SettingsRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Setting, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.Setting, error)
	Create(ctx context.Context, setting models.Setting) error
	Update(ctx context.Context, setting models.Setting) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
}

type settingsRepository struct {
	client repositories.PostgresClient
}

func New(client repositories.PostgresClient) SettingsRepositorer {
	return &settingsRepository{
		client: client,
	}
}

func (r *settingsRepository) FindByID(ctx context.Context, id int) (models.Setting, error) {
	var (
		query = `
			SELECT s.id, s.name, c.id, c.name
			FROM settings s
			JOIN settings_categories c ON c.id = s.id_category
			WHERE s.id = $1
			LIMIT 1
		`
		setting models.Setting
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&setting.ID, &setting.Name, &setting.Category.ID, &setting.Category.Name); err != nil {
		return models.Setting{}, fmt.Errorf("failed to find setting by id: %w", err)
	}
	return setting, nil
}

func (r *settingsRepository) FindMany(ctx context.Context, limit, offset int, name string) ([]models.Setting, error) {
	var (
		query = `
			SELECT s.id, s.name, c.id, c.name
			FROM settings s
			JOIN settings_categories c ON c.id = s.id_category
			WHERE s.name LIKE $1
			ORDER BY c.id
			LIMIT $2
			OFFSET $3
		`
		settings []models.Setting
	)
	rows, err := r.client.Query(ctx, query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find settings: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		setting := models.Setting{}
		if err := rows.Scan(&setting.ID, &setting.Name, &setting.Category.ID, &setting.Category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func (r *settingsRepository) Create(ctx context.Context, setting models.Setting) error {
	var (
		query = `
			INSERT INTO settings
				(name, id_category)
			VALUES
				($1, $2)
		`
	)
	if _, err := r.client.Exec(ctx, query, setting.Name, setting.Category.ID); err != nil {
		return fmt.Errorf("failed to create setting: %w", err)
	}
	return nil
}

func (r *settingsRepository) Update(ctx context.Context, setting models.Setting) error {
	var (
		query = `
			UPDATE settings
			SET name = $1,
				id_category = $2,
				updated_at = NOW()
			WHERE id = $3
		`
	)
	if _, err := r.client.Exec(ctx, query, setting.Name, setting.Category.ID, setting.ID); err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}
	return nil
}

func (r *settingsRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE settings
			SET deleted_at = NOW(),
				updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}
	return nil
}

func (r *settingsRepository) Restore(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE settings
			SET deleted_at = NULL,
				updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}
	return nil
}
