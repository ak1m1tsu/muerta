package setting

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type SettingsRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Setting, error)
	FindMany(ctx context.Context, filter models.SettingFilter) ([]models.Setting, error)
	Create(ctx context.Context, setting models.Setting) error
	Update(ctx context.Context, setting models.Setting) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) (models.Setting, error)
	Count(ctx context.Context, filter models.SettingFilter) (int, error)
}

type settingsRepository struct {
	client repositories.PostgresClient
}

func New(client repositories.PostgresClient) SettingsRepositorer {
	return &settingsRepository{
		client: client,
	}
}

func (r *settingsRepository) Count(ctx context.Context, filter models.SettingFilter) (int, error) {
	var (
		query = `
			SELECT COUNT(*) 
			FROM settings 
			WHERE deleted_at IS NULL AND
				name LIKE $1
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query, "%"+filter.Name+"%").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count settings: %w", err)
	}
	return count, nil
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

func (r *settingsRepository) FindMany(
	ctx context.Context,
	filter models.SettingFilter,
) ([]models.Setting, error) {
	var (
		query = `
			SELECT s.id, s.name, c.id, c.name
			FROM settings s
			JOIN settings_categories c ON c.id = s.id_category
			WHERE 
				s.name ILIKE $1 AND
				s.deleted_at IS NULL
			ORDER BY c.id
			LIMIT $2
			OFFSET $3
		`
		settings []models.Setting
	)
	rows, err := r.client.Query(ctx, query, "%"+filter.Name+"%", filter.Limit, filter.Offset)
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
	query := `
			INSERT INTO settings
				(name, id_category)
			VALUES
				($1, $2)
		`
	if _, err := r.client.Exec(ctx, query, setting.Name, setting.Category.ID); err != nil {
		return fmt.Errorf("failed to create setting: %w", err)
	}
	return nil
}

func (r *settingsRepository) Update(ctx context.Context, setting models.Setting) error {
	query := `
			UPDATE settings
			SET name = $1,
				id_category = $2,
				updated_at = NOW()
			WHERE id = $3
		`
	if _, err := r.client.Exec(ctx, query, setting.Name, setting.Category.ID, setting.ID); err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}
	return nil
}

func (r *settingsRepository) Delete(ctx context.Context, id int) error {
	query := `
			UPDATE settings
			SET deleted_at = NOW(),
				updated_at = NOW()
			WHERE id = $1
		`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}
	return nil
}

func (r *settingsRepository) Restore(ctx context.Context, id int) (models.Setting, error) {
	var (
		query = `
			WITH updated AS (
				UPDATE settings
				SET deleted_at = NULL,
					updated_at = NOW()
				WHERE id = $1
				RETURNING id, name, id_category
			)
			SELECT u.id, u.name, u.id_category, c.name
			FROM updated u
			JOIN settings_categories c ON c.id = u.id_category
			WHERE u.id_category = c.id
		`
		model models.Setting
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&model.ID, &model.Name, &model.Category.ID, &model.Category.Name); err != nil {
		return models.Setting{}, fmt.Errorf("failed to restore setting: %w", err)
	}
	return model, nil
}
