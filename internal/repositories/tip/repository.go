package tip

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type TipRepositorer interface {
	repositories.Repository
	FindByID(ctx context.Context, id int) (models.Tip, error)
	FindMany(ctx context.Context, limit, offset int, description string) ([]models.Tip, error)
	Create(ctx context.Context, tip *models.Tip) error
	Update(ctx context.Context, tip models.Tip) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	FindProducts(ctx context.Context, id int) ([]models.Product, error)
	FindStorages(ctx context.Context, id int) ([]models.Storage, error)
	AddProduct(ctx context.Context, tipID, productID int) (models.Product, error)
	RemoveProduct(ctx context.Context, tipID, productID int) error
	AddStorage(ctx context.Context, tipID, storageID int) (models.Storage, error)
	RemoveStorage(ctx context.Context, tipID, storageID int) error
}

type tipRepository struct {
	client repositories.PostgresClient
}

// AddProduct implements TipRepositorer
func (r *tipRepository) AddProduct(ctx context.Context, tipID int, productID int) (models.Product, error) {
	var (
		query = `
			WITH inserted AS (
				INSERT INTO products_tips (id_tip, id_product)
				VALUES ($1, $2)
				RETURNING id_tip, id_product
			)
			SELECT p.id, p.name
			FROM products p
			JOIN inserted i ON i.id_product = p.id
			WHERE p.id = i.id_product AND p.deleted_at IS NULL
		`
		model models.Product
	)
	if err := r.client.QueryRow(ctx, query, tipID, productID).Scan(&model.ID, &model.Name); err != nil {
		return models.Product{}, fmt.Errorf("failed to add product: %w", err)
	}
	return model, nil
}

// AddStoragep implements TipRepositorer
func (r *tipRepository) AddStorage(ctx context.Context, tipID int, storageID int) (models.Storage, error) {
	var (
		query = `
			WITH inserted AS (
				INSERT INTO storages_tips (id_tip, id_storage)
				VALUES ($1, $2)
				RETURNING id_tip, id_storage
			)
			SELECT s.id, s.name, st.id, st.name, s.temperature, s.humidity
			FROM storages s
			JOIN storages_types st ON s.id_type = st.id
			JOIN inserted i ON i.id_storage = s.id
			WHERE s.id = i.id_storage AND s.deleted_at IS NULL
		`
		model models.Storage
	)
	if err := r.client.QueryRow(ctx, query, tipID, storageID).Scan(
		&model.ID, &model.Name, &model.Type.ID, &model.Type.Name, &model.Temperature, &model.Humidity,
	); err != nil {
		return models.Storage{}, fmt.Errorf("failed to add storage: %w", err)
	}
	return model, nil
}

// RemoveProduct implements TipRepositorer
func (r *tipRepository) RemoveProduct(ctx context.Context, tipID int, productID int) error {
	var query = `
		DELETE FROM products_tips
		WHERE id_tip = $1 AND id_product = $2
	`
	if _, err := r.client.Exec(ctx, query, tipID, productID); err != nil {
		return fmt.Errorf("failed to remove products: %w", err)
	}
	return nil
}

// RemoveStorage implements TipRepositorer
func (r *tipRepository) RemoveStorage(ctx context.Context, tipID int, storageID int) error {
	var query = `
		DELETE FROM storages_tips
		WHERE id_tip = $1 AND id_storage = $2
	`
	if _, err := r.client.Exec(ctx, query, tipID, storageID); err != nil {
		return fmt.Errorf("failed to remove storage: %w", err)
	}
	return nil
}

func New(client repositories.PostgresClient) TipRepositorer {
	return &tipRepository{
		client: client,
	}
}

func (r *tipRepository) Count(ctx context.Context) (int, error) {
	var (
		query = `
			SELECT COUNT(*) FROM tips WHERE deleted_at IS NULL
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count tips: %w", err)
	}
	return count, nil
}

// FindProducts implements TipRepositorer
func (r *tipRepository) FindProducts(ctx context.Context, id int) ([]models.Product, error) {
	var (
		query = `
			SELECT p.id, p.name
			FROM products p
			JOIN products_tips pt ON p.id = pt.id_product
			WHERE pt.id_tip = $1 AND p.deleted_at IS NULL
		`
		products []models.Product
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return products, fmt.Errorf("failed to find products: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name); err != nil {
			return products, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}

// FindStorages implements TipRepositorer
func (r *tipRepository) FindStorages(ctx context.Context, id int) ([]models.Storage, error) {
	var (
		query = `
			SELECT s.id, s.name
			FROM storages s
			JOIN storages_tips st ON s.id = st.id_storage
			WHERE st.id_tip = $1 AND s.deleted_at IS NULL
		`
		storages []models.Storage
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return storages, fmt.Errorf("failed to find storages: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var storage models.Storage
		if err := rows.Scan(&storage.ID, &storage.Name); err != nil {
			return storages, fmt.Errorf("failed to scan storage: %w", err)
		}
		storages = append(storages, storage)
	}
	return storages, nil
}

// Create implements TipRepositorer
func (r *tipRepository) Create(ctx context.Context, tip *models.Tip) error {
	var (
		query = `
			INSERT INTO tips (description)
			VALUES ($1)
			RETURNING id
		`
	)
	if err := r.client.QueryRow(ctx, query, tip.Description).Scan(&tip.ID); err != nil {
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
