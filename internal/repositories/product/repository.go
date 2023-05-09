package product

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type ProductRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Product, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.Product, error)
	Create(ctx context.Context, product models.Product) error
	Update(ctx context.Context, product models.Product) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	FindCategories(ctx context.Context, id int) ([]models.ProductCategory, error)
	FindRecipes(ctx context.Context, id int) ([]models.Recipe, error)
}

type productRepository struct {
	client repositories.PostgresClient
}

func New(client repositories.PostgresClient) ProductRepositorer {
	return &productRepository{client: client}
}

// FindCategories implements ProductRepositorer
func (r *productRepository) FindCategories(ctx context.Context, id int) ([]models.ProductCategory, error) {
	var (
		query = `
			SELECT c.id, c.name
			FROM categories c
			JOIN products_categories pc ON pc.id_category = c.id
			WHERE pc.id_product = $1 AND c.deleted_at IS NULL
		`
		categories []models.ProductCategory
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var category models.ProductCategory
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// FindRecipes implements ProductRepositorer
func (r *productRepository) FindRecipes(ctx context.Context, id int) ([]models.Recipe, error) {
	var (
		query = `
			SELECT r.id, r.name
			FROM recipes r
			JOIN products_recipes_measures prm ON prm.id_recipe = r.id
			WHERE prm.id_product = $1 AND r.deleted_at IS NULL
		`
		recipes []models.Recipe
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var category models.Recipe
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		recipes = append(recipes, category)
	}
	return recipes, nil
}

func (repo *productRepository) FindByID(ctx context.Context, id int) (models.Product, error) {
	var (
		query = `
			SELECT id, name
			FROM products
			WHERE id = $1
			LIMIT 1
		`
		product models.Product
	)
	if err := repo.client.QueryRow(ctx, query, id).Scan(&product.ID, &product.Name); err != nil {
		return models.Product{}, fmt.Errorf("failed to find product: %w", err)
	}
	return product, nil
}

func (repo *productRepository) FindMany(ctx context.Context, limit, offset int, name string) ([]models.Product, error) {
	var (
		query = `
			SELECT id, name
			FROM products
			WHERE name ILIKE $3
			ORDER BY name ASC
			LIMIT $1
			OFFSET $2
		`
		products = make([]models.Product, 0, limit)
	)
	rows, err := repo.client.Query(ctx, query, limit, offset, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}

func (repo *productRepository) Create(ctx context.Context, product models.Product) error {
	var (
		query = `
			INSERT INTO products (name)
			VALUES ($1)
		`
	)
	if _, err := repo.client.Exec(ctx, query, product.Name); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (repo *productRepository) Update(ctx context.Context, product models.Product) error {
	var (
		query = `
			UPDATE products
			SET name = $1,
				updated_at = NOW()
			WHERE id = $2
		`
	)
	if _, err := repo.client.Exec(ctx, query, product.Name, product.ID); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (repo *productRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE products
			SET deleted_at = NOW(),
				updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := repo.client.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (repo *productRepository) Restore(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE products
			SET deleted_at = NULL,
				updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := repo.client.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}
