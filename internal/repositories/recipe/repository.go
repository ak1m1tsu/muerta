package recipes

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type RecipeIngredientsRepositorer interface {
	FindIngredients(ctx context.Context, id int) ([]models.RecipeIngredient, error)
	CreateIngredient(ctx context.Context, id int, entity *models.RecipeIngredient) (models.RecipeIngredient, error)
	UpdateIngredient(ctx context.Context, id int, entity *models.RecipeIngredient) (models.RecipeIngredient, error)
	DeleteIngredient(ctx context.Context, recipeId, productId int) error
}

type RecipeStepsRepositorer interface {
	FindSteps(ctx context.Context, recipeID int) ([]models.Step, error)
	CreateStep(ctx context.Context, recipeID, stepID, place int) (models.Step, error)
	DeleteStep(ctx context.Context, recipeID, stepID, place int) error
}

type RecipesRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Recipe, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.Recipe, error)
	Create(ctx context.Context, recipe *models.Recipe) error
	Update(ctx context.Context, recipe *models.Recipe) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	RecipeIngredientsRepositorer
	RecipeStepsRepositorer
	repositories.Repository
}

type recipesRepository struct {
	client repositories.PostgresClient
}

func (r *recipesRepository) Count(ctx context.Context) (int, error) {
	var (
		query = `
			SELECT COUNT(*) FROM recipes WHERE deleted_at IS NULL
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count recipes: %w", err)
	}
	return count, nil
}

// CreateStep implements RecipesRepositorer
func (r *recipesRepository) CreateStep(ctx context.Context, recipeID, stepID, place int) (models.Step, error) {
	var (
		query = `
			WITH inserted AS (
				INSERT INTO recipes_steps (id_recipe, id_step, place) 
				VALUES ($1, $2, $3)
				RETURNING id_recipe, id_step, place
			)
			SELECT s.id, s.name, i.place
			FROM steps s
			JOIN inserted i ON i.id_step = s.id
			WHERE s.id = $2
			LIMIT 1
		`
		result models.Step
	)
	if err := r.client.QueryRow(ctx, query, recipeID, stepID, place).Scan(&result.ID, &result.Name, &result.Place); err != nil {
		return models.Step{}, fmt.Errorf("failed to create step: %w", err)
	}
	return result, nil
}

// DeleteStep implements RecipesRepositorer
func (r *recipesRepository) DeleteStep(ctx context.Context, recipeID, stepID, place int) error {
	query := `
		DELETE FROM recipes_steps 
		WHERE id_recipe = $1 AND 
			id_step = $2 AND 
			place = $3
	`
	if _, err := r.client.Exec(ctx, query, recipeID, stepID, place); err != nil {
		return fmt.Errorf("failed to delete step: %w", err)
	}
	return nil
}

// FindSteps implements RecipesRepositorer
func (r *recipesRepository) FindSteps(ctx context.Context, recipeID int) ([]models.Step, error) {
	var (
		query = `
			SELECT s.id, s.name, rs.place
			FROM steps s
			JOIN recipes_steps rs ON rs.id_step = s.id
			WHERE rs.id_recipe = $1
		`
		result []models.Step
	)
	rows, err := r.client.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find steps: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var step models.Step
		if err := rows.Scan(&step.ID, &step.Name, &step.Place); err != nil {
			return nil, fmt.Errorf("failed to scan step: %w", err)
		}
		result = append(result, step)
	}
	return result, nil
}

// CreateIngredient implements RecipesRepositorer
func (r *recipesRepository) CreateIngredient(ctx context.Context, id int, entity *models.RecipeIngredient) (models.RecipeIngredient, error) {
	var (
		queryInsert = `
			INSERT INTO products_recipes_measures 
				(id_product, id_recipe, id_measure, quantity)
			VALUES
				($1, $2, $3, $4)
		`
		querySelect = `
			SELECT p.name, m.name
			FROM products_recipes_measures prm
			JOIN products p ON p.id = prm.id_product
			JOIN measures m ON m.id = prm.id_measure
			WHERE prm.id_recipe = $1 AND 
				  prm.id_product = $2 AND 
				  prm.id_measure = $3
			LIMIT 1
		`
	)
	if _, err := r.client.Exec(ctx, queryInsert, entity.Product.ID, id, entity.Measure.ID, entity.Quantity); err != nil {
		return models.RecipeIngredient{}, fmt.Errorf("failed to create recipe ingredient: %w", err)
	}
	if err := r.client.QueryRow(ctx, querySelect, id, entity.Product.ID, entity.Measure.ID).Scan(&entity.Product.Name, &entity.Measure.Name); err != nil {
		return models.RecipeIngredient{}, fmt.Errorf("failed to query recipe ingredient: %w", err)
	}
	return *entity, nil
}

// DeleteIngredient implements RecipesRepositorer
func (r *recipesRepository) DeleteIngredient(ctx context.Context, recipeId int, productId int) error {
	var (
		query = `
			DELETE FROM products_recipes_measures
			WHERE id_recipe = $1 AND id_product = $2
		`
	)
	if _, err := r.client.Exec(ctx, query, recipeId, productId); err != nil {
		return fmt.Errorf("failed to delete recipe ingredient: %w", err)
	}
	return nil
}

// FindIngredients implements RecipesRepositorer
func (r *recipesRepository) FindIngredients(ctx context.Context, id int) ([]models.RecipeIngredient, error) {
	var (
		query = `
		SELECT p.id, p.name, m.id, m.name, prm.quantity
			FROM products_recipes_measures prm
			JOIN products p ON p.id = prm.id_product
			JOIN measures m ON m.id = prm.id_measure
		WHERE prm.id_recipe = $1
		`
		entities []models.RecipeIngredient
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query ingredients: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var entity models.RecipeIngredient
		if err := rows.Scan(&entity.Product.ID, &entity.Product.Name, &entity.Measure.ID, &entity.Measure.Name, &entity.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

// UpdateIngredient implements RecipesRepositorer
func (r *recipesRepository) UpdateIngredient(ctx context.Context, id int, entity *models.RecipeIngredient) (models.RecipeIngredient, error) {
	var (
		query = `
			UPDATE products_recipes_measures
			SET id_measure = $2, 
				quantity = $3
			WHERE id_recipe = $4 AND 
				  id_product = $1
		`
		querySelect = `
			SELECT m.name, prm.quantity
			FROM products_recipes_measures prm
			JOIN products p ON p.id = prm.id_product
			JOIN measures m ON m.id = prm.id_measure
			WHERE prm.id_recipe = $1 AND 
				  prm.id_product = $2 AND 
				  prm.id_measure = $3
			LIMIT 1
		`
	)
	if _, err := r.client.Exec(ctx, query, entity.Product.ID, entity.Measure.ID, entity.Quantity, id); err != nil {
		return models.RecipeIngredient{}, fmt.Errorf("failed to update recipe ingredient: %w", err)
	}
	if err := r.client.QueryRow(ctx, querySelect, id, entity.Product.ID, entity.Measure.ID).Scan(&entity.Measure.Name, &entity.Quantity); err != nil {
		return models.RecipeIngredient{}, fmt.Errorf("failed to query recipe ingredient: %w", err)
	}
	return *entity, nil
}

func New(client repositories.PostgresClient) RecipesRepositorer {
	return &recipesRepository{client: client}
}

func (r *recipesRepository) FindByID(ctx context.Context, id int) (models.Recipe, error) {
	var (
		query = `
			SELECT id, name, description 
			FROM recipes 
			WHERE id = $1
		`
		querySteps = `
			SELECT rs.id_step, rs.place, s.name
			FROM recipes_steps rs
			JOIN steps s ON s.id = rs.id_step
			WHERE rs.id_recipe = $1
			ORDER BY rs.place ASC
		`
		recipe models.Recipe
	)
	if err := r.client.QueryRow(ctx, query, id).Scan(&recipe.ID, &recipe.Name, &recipe.Description); err != nil {
		return models.Recipe{}, fmt.Errorf("failed to query recipe: %w", err)
	}
	rows, err := r.client.Query(ctx, querySteps, id)
	if err != nil {
		return models.Recipe{}, fmt.Errorf("failed to query steps: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var step models.Step
		if err := rows.Scan(&step.ID, &step.Place, &step.Name); err != nil {
			return models.Recipe{}, fmt.Errorf("failed to scan step: %w", err)
		}
		recipe.Steps = append(recipe.Steps, step)
	}
	return recipe, nil
}

func (r *recipesRepository) FindMany(ctx context.Context, limit, offset int, name string) ([]models.Recipe, error) {
	var (
		query = `
			SELECT id, name, description
			FROM recipes
			WHERE 
				name LIKE $1 AND 
				deleted_at IS NULL
			ORDER BY id ASC
			LIMIT $2
			OFFSET $3
		`
		recipes []models.Recipe
	)
	rows, err := r.client.Query(ctx, query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.Description); err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *recipesRepository) Create(ctx context.Context, recipe *models.Recipe) error {
	var (
		query = `
			INSERT INTO recipes
				(id_user, name, description)
			VALUES
				($1, $2, $3)
			RETURNING id
		`
	)
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, query, recipe.User.ID, recipe.Name, recipe.Description).Scan(&recipe.ID); err != nil {
		return fmt.Errorf("failed to create recipe: %w", err)
	}
	if _, err = tx.CopyFrom(ctx,
		pgx.Identifier{"recipes_steps"},
		[]string{"id_recipe", "id_step", "place"},
		pgx.CopyFromSlice(len(recipe.Steps), func(i int) ([]any, error) {
			return []any{recipe.ID, recipe.Steps[i].ID, recipe.Steps[i].Place}, nil
		}),
	); err != nil {
		return err
	}
	if _, err = tx.CopyFrom(ctx,
		pgx.Identifier{"products_recipes_measures"},
		[]string{"id_product", "id_recipe", "id_measure", "quantity"},
		pgx.CopyFromSlice(len(recipe.Ingredients), func(i int) ([]any, error) {
			return []any{recipe.Ingredients[i].Product.ID, recipe.ID, recipe.Ingredients[i].Measure.ID, recipe.Ingredients[i].Quantity}, nil
		}),
	); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *recipesRepository) Update(ctx context.Context, recipe *models.Recipe) error {
	var (
		query = `
			UPDATE recipes
			SET name = $1,
				description = $2,
				updated_at = NOW()
			WHERE id = $3
		`
	)
	if _, err := r.client.Exec(ctx, query, recipe.Name, recipe.Description, recipe.ID); err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}
	return nil
}

func (r *recipesRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE recipes
			SET deleted_at = NOW(),
				updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}
	return nil
}

func (r *recipesRepository) Restore(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE recipes
			SET deleted_at = NULL,
				updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := r.client.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("failed to restore recipe: %w", err)
	}
	return nil
}
