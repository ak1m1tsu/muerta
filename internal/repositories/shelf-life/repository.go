package shelflife

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type ShelfLifeRepositorer interface {
	FindByID(ctx context.Context, id int) (models.ShelfLife, error)
	FindMany(ctx context.Context, filter models.ShelfLifeFilter) ([]models.ShelfLife, error)
	Create(ctx context.Context, measure models.ShelfLife) error
	Update(ctx context.Context, measure models.ShelfLife) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	CreateStatus(ctx context.Context, id, statusID int) (models.ShelfLifeStatus, error)
	DeleteStatus(ctx context.Context, id, statusID int) error
	FindStatuses(ctx context.Context, id int) ([]models.ShelfLifeStatus, error)
	Count(ctx context.Context, filter models.ShelfLifeFilter) (int, error)
}

type shelfLifeRepository struct {
	client repositories.PostgresClient
}

func (r *shelfLifeRepository) Count(
	ctx context.Context,
	filter models.ShelfLifeFilter,
) (int, error) {
	var (
		query = `
			SELECT COUNT(*) FROM shelf_lives WHERE deleted_at IS NULL
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count shelf lives: %w", err)
	}
	return count, nil
}

// CreateStatus implements ShelfLifeRepositorer
func (r *shelfLifeRepository) CreateStatus(
	ctx context.Context,
	id int,
	statusID int,
) (models.ShelfLifeStatus, error) {
	var (
		query = `
			WITH inserted AS (
				INSERT INTO shelf_lives_statuses (id_shelf_life, id_status)
				VALUES ($1, $2)
				RETURNING id_shelf_life, id_status
			)
			SELECT s.id, s.name
			FROM statuses s
			JOIN inserted i ON i.id_status = s.id
			WHERE s.id = $2
			LIMIT 1
		`
		result models.ShelfLifeStatus
	)
	if err := r.client.QueryRow(ctx, query, id, statusID).Scan(&result.ID, &result.Name); err != nil {
		return models.ShelfLifeStatus{}, fmt.Errorf("failed to create shelf life status: %w", err)
	}
	return result, nil
}

// DeleteStatus implements ShelfLifeRepositorer
func (r *shelfLifeRepository) DeleteStatus(ctx context.Context, id int, statusID int) error {
	query := `
		DELETE FROM shelf_lives_statuses
		WHERE id_shelf_life = $1 AND id_status = $2
	`
	if _, err := r.client.Exec(ctx, query, id, statusID); err != nil {
		return fmt.Errorf("failed to delete shelf life status: %w", err)
	}
	return nil
}

// FindStatuses implements ShelfLifeRepositorer
func (r *shelfLifeRepository) FindStatuses(
	ctx context.Context,
	id int,
) ([]models.ShelfLifeStatus, error) {
	var (
		query = `
			SELECT s.id, s.name
			FROM statuses s
			JOIN shelf_lives_statuses sl ON sl.id_status = s.id
			WHERE sl.id_shelf_life = $1
		`
		result []models.ShelfLifeStatus
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find shelf life statuses: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var status models.ShelfLifeStatus
		if err := rows.Scan(&status.ID, &status.Name); err != nil {
			return nil, fmt.Errorf("failed to scan shelf life status: %w", err)
		}
		result = append(result, status)
	}
	return result, nil
}

// Create implements ShelfLifeRepositorer
func (r *shelfLifeRepository) Create(ctx context.Context, model models.ShelfLife) error {
	query := `
			INSERT INTO shelf_lives
				(id, id_product, id_storage, id_measure, id_user, quantity, purchase_date, end_date)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8)
		`
	if _, err := r.client.Exec(ctx, query, model.ID, model.Product.ID, model.Storage.ID, model.Measure.ID, model.User.ID, model.Quantity, model.PurchaseDate, model.EndDate); err != nil {
		return fmt.Errorf("failed to create shelf life: %w", err)
	}
	return nil
}

// Delete implements ShelfLifeRepositorer
func (r *shelfLifeRepository) Delete(ctx context.Context, id int) error {
	query := `
			UPDATE shelf_lives
			SET deleted_at = NOW(),
				updated_at = NOW()
			WHERE id = $1
		`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete shelf life: %w", err)
	}
	return nil
}

// FindByID implements ShelfLifeRepositorer
func (r *shelfLifeRepository) FindByID(ctx context.Context, id int) (models.ShelfLife, error) {
	var (
		query = `
			SELECT 
				sl.id, 
				sl.id_product, p.name,
				sl.id_storage, s.name,
				sl.id_measure, m.name,
				sl.quantity, sl.purchase_date, sl.end_date
			FROM shelf_lives sl
			JOIN products p ON p.id = sl.id_product
			JOIN storages s ON s.id = sl.id_storage
			JOIN measures m ON m.id = sl.id_measure
			WHERE sl.id = $1
			LIMIT 1
		`
		model models.ShelfLife
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&model.ID, &model.Product.ID, &model.Product.Name, &model.Storage.ID, &model.Storage.Name, &model.Measure.ID, &model.Measure.Name, &model.Quantity, &model.PurchaseDate, &model.EndDate); err != nil {
		return models.ShelfLife{}, fmt.Errorf("failed to find shelf life: %w", err)
	}
	return model, nil
}

// FindMany implements ShelfLifeRepositorer
func (r *shelfLifeRepository) FindMany(
	ctx context.Context,
	filter models.ShelfLifeFilter,
) ([]models.ShelfLife, error) {
	var (
		query = `
			SELECT sl.id, 
				sl.id_product, p.name,
				sl.id_storage, s.name,
				sl.id_measure, m.name,
				sl.quantity, sl.purchase_date, sl.end_date
			FROM shelf_lives sl
			JOIN products p ON p.id = sl.id_product
			JOIN storages s ON s.id = sl.id_storage
			JOIN measures m ON m.id = sl.id_measure
			WHERE sl.deleted_at IS NULL
			ORDER BY sl.created_at DESC
			LIMIT $1 
			OFFSET $2
		`
		shelfLives = make([]models.ShelfLife, 0, filter.Limit)
	)
	rows, err := r.client.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find shelf lives: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var shelfLife models.ShelfLife
		if err := rows.Scan(&shelfLife.ID, &shelfLife.Product.ID, &shelfLife.Product.Name, &shelfLife.Storage.ID, &shelfLife.Storage.Name, &shelfLife.Measure.ID, &shelfLife.Measure.Name, &shelfLife.Quantity, &shelfLife.PurchaseDate, &shelfLife.EndDate); err != nil {
			return nil, fmt.Errorf("failed to scan shelf life: %w", err)
		}
		shelfLives = append(shelfLives, shelfLife)
	}
	return shelfLives, nil
}

// Restore implements ShelfLifeRepositorer
func (r *shelfLifeRepository) Restore(ctx context.Context, id int) error {
	query := `
			UPDATE shelf_lives
			SET deleted_at = NULL,
				updated_at = NOW()
			WHERE id = $1
		`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to restore shelf life: %w", err)
	}
	return nil
}

// Update implements ShelfLifeRepositorer
func (r *shelfLifeRepository) Update(ctx context.Context, model models.ShelfLife) error {
	query := `
			UPDATE shelf_lives
			SET id_product = $1,
				id_storage = $2,
				id_measure = $3,
				quantity = $4,
				purchase_date = $5,
				end_date = $6,
				updated_at = NOW()
			WHERE id = $7
		`
	if _, err := r.client.Exec(ctx, query, model.Product.ID, model.Storage.ID, model.Measure.ID, model.Quantity, model.PurchaseDate, model.EndDate, model.ID); err != nil {
		return fmt.Errorf("failed to update shelf life: %w", err)
	}
	return nil
}

func New(client repositories.PostgresClient) ShelfLifeRepositorer {
	return &shelfLifeRepository{
		client: client,
	}
}
