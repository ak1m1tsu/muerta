package storage

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

type StorageRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Vault, error)
	FindMany(ctx context.Context, filter models.StorageFilter) ([]models.Vault, error)
	Create(ctx context.Context, storage *models.Vault) error
	Update(ctx context.Context, storage *models.Vault) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	CreateTip(ctx context.Context, id, tipID int) (models.Tip, error)
	DeleteTip(ctx context.Context, id, tipID int) error
	FindTips(ctx context.Context, id int) ([]models.Tip, error)
	FindShelfLives(ctx context.Context, id int) ([]models.ShelfLife, error)
	Count(ctx context.Context, filter models.StorageFilter) (int, error)
}

type storageRepository struct {
	client postgres.Client
}

func (r *storageRepository) FindShelfLives(
	ctx context.Context,
	id int,
) ([]models.ShelfLife, error) {
	var (
		query = `
			SELECT sl.id, p.id, p.name, m.id, m.name, sl.quantity, sl.purchase_date, sl.end_date
			FROM shelf_lives sl
			JOIN products p ON p.id = sl.id_product
			JOIN measures m ON m.id = sl.id_measure
			WHERE sl.deleted_at IS NULL AND sl.id_storage = $1
			ORDER BY sl.end_date DESC, sl.purchase_date DESC
		`
		result []models.ShelfLife
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("find shelf lives: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var model models.ShelfLife
		if err := rows.Scan(
			&model.ID,
			&model.Product.ID,
			&model.Product.Name,
			&model.Measure.ID,
			&model.Measure.Name,
			&model.Quantity,
			&model.PurchaseDate,
			&model.EndDate,
		); err != nil {
			return nil, fmt.Errorf("find shelf lives: %w", err)
		}
		result = append(result, model)
	}
	return result, nil
}

func (r *storageRepository) Count(ctx context.Context, filter models.StorageFilter) (int, error) {
	var (
		query = `
			SELECT COUNT(*) 
			FROM storages 
			WHERE deleted_at IS NULL AND
				name ILIKE $1
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query, "%"+filter.Name+"%").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count storages: %w", err)
	}
	return count, nil
}

// CreateTip implements StorageRepositorer
func (r *storageRepository) CreateTip(ctx context.Context, id int, tipID int) (models.Tip, error) {
	var (
		query = `
			WITH inserted AS (
				INSERT INTO storages_tips (id_storage, id_tip)
				VALUES ($1, $2)
				RETURNING id_storage, id_tip
			)
			SELECT t.id, t.description
			FROM tips t
			JOIN inserted i ON i.id_tip = t.id
			WHERE t.id = $2 AND t.deleted_at IS NULL
			LIMIT 1
		`
		tip models.Tip
	)
	if err := r.client.QueryRow(ctx, query, id, tipID).Scan(&tip.ID, &tip.Description); err != nil {
		return models.Tip{}, fmt.Errorf("create tip: %w", err)
	}
	return tip, nil
}

// DeleteTip implements StorageRepositorer
func (r *storageRepository) DeleteTip(ctx context.Context, id int, tipID int) error {
	query := `
			DELETE FROM storages_tips
			WHERE id_storage = $1 AND id_tip = $2
		`
	if _, err := r.client.Exec(ctx, query, id, tipID); err != nil {
		return fmt.Errorf("delete tip: %w", err)
	}
	return nil
}

// FindTips implements StorageRepositorer
func (r *storageRepository) FindTips(ctx context.Context, id int) ([]models.Tip, error) {
	var (
		query = `
			SELECT t.id, t.description
			FROM tips t
			JOIN storages_tips st ON st.id_tip = t.id
			WHERE st.id_storage = $1 AND t.deleted_at IS NULL
		`
		tips []models.Tip
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("find tips: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var tip models.Tip
		if err := rows.Scan(&tip.ID, &tip.Description); err != nil {
			return nil, fmt.Errorf("find tips: %w", err)
		}
		tips = append(tips, tip)
	}
	return tips, nil
}

func New(client postgres.Client) StorageRepositorer {
	return &storageRepository{
		client: client,
	}
}

func (r *storageRepository) FindByID(ctx context.Context, id int) (models.Vault, error) {
	var (
		query = `
			SELECT 
				s.id, s.name, 
				s.temperature, s.humidity,
				s.created_at, st.id, st.name
			FROM storages s
			JOIN storages_types st ON st.id = s.id_type
			WHERE s.id = $1 AND s.deleted_at IS NULL
			LIMIT 1
		`
		storage models.Vault
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&storage.ID, &storage.Name, &storage.Temperature, &storage.Humidity, &storage.CreatedAt, &storage.Type.ID, &storage.Type.Name); err != nil {
		return models.Vault{}, fmt.Errorf("find storage by id: %w", err)
	}
	return storage, nil
}

func (r *storageRepository) FindMany(
	ctx context.Context,
	filter models.StorageFilter,
) ([]models.Vault, error) {
	var (
		query = `
			SELECT 
				s.id, s.name, 
				s.temperature, s.humidity,
				s.created_at, st.id, st.name
			FROM storages s
			JOIN storages_types st ON st.id = s.id_type
			WHERE s.deleted_at IS NULL
				AND s.name ILIKE $3
			ORDER BY s.created_at DESC
			LIMIT $1 OFFSET $2
		`
		storages []models.Vault
	)
	rows, err := r.client.Query(ctx, query, filter.Limit, filter.Offset, "%"+filter.Name+"%")
	if err != nil {
		return nil, fmt.Errorf("find many storages: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var storage models.Vault
		if err := rows.Scan(&storage.ID, &storage.Name, &storage.Temperature, &storage.Humidity, &storage.CreatedAt, &storage.Type.ID, &storage.Type.Name); err != nil {
			return nil, fmt.Errorf("scan storage: %w", err)
		}
		storages = append(storages, storage)
	}
	return storages, nil
}

func (r *storageRepository) Create(ctx context.Context, storage *models.Vault) error {
	query := `
			INSERT INTO storages 
				(name, temperature, humidity, id_type)
			VALUES
				($1, $2, $3, $4)
			RETURNING id
		`
	if _, err := r.client.Exec(ctx, query, storage.Name, storage.Temperature, storage.Humidity, storage.Type.ID); err != nil {
		return fmt.Errorf("create storage: %w", err)
	}
	return nil
}

func (r *storageRepository) Update(ctx context.Context, storage *models.Vault) error {
	query := `
			UPDATE storages
			SET name = $1,
				temperature = $2,
				humidity = $3,
				id_type = $4,
				updated_at = NOW()
			WHERE id = $5
		`
	if _, err := r.client.Exec(ctx, query, storage.Name, storage.Temperature, storage.Humidity, storage.Type.ID, storage.ID); err != nil {
		return fmt.Errorf("update storage: %w", err)
	}
	return nil
}

func (r *storageRepository) Delete(ctx context.Context, id int) error {
	query := `
			UPDATE storages
			SET deleted_at = NOW(),
				updated_at = NOW()
			WHERE id = $1
		`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("delete storage: %w", err)
	}
	return nil
}

func (r *storageRepository) Restore(ctx context.Context, id int) error {
	query := `
			UPDATE storages
			SET deleted_at = NULL,
				updated_at = NOW()
			WHERE id = $1
		`
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("restore storage: %w", err)
	}
	return nil
}
