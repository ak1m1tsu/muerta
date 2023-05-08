package recipes

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type RecipesRepositorer interface {
	FindByID(ctx context.Context, id int) (models.Recipe, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.Recipe, error)
	Create(ctx context.Context, recipe *models.Recipe) error
	Update(ctx context.Context, recipe *models.Recipe) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
}

type RecipesRepository struct {
	client repositories.PostgresClient
}

func New(client repositories.PostgresClient) *RecipesRepository {
	return &RecipesRepository{client: client}
}

func (r *RecipesRepository) FindByID(ctx context.Context, id int) (models.Recipe, error) {
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

func (r *RecipesRepository) FindMany(ctx context.Context, limit, offset int, name string) ([]models.Recipe, error) {
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

func (r *RecipesRepository) Create(ctx context.Context, recipe *models.Recipe) error {
	var (
		query = `
			INSERT INTO recipes
				(name, description)
			VALUES
				($1, $2)
			RETURNING id
		`
	)
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, query, recipe.Name, recipe.Description).Scan(&recipe.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			return fmt.Errorf(fmt.Sprintf(
				"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message,
				pgErr.Detail,
				pgErr.Where,
				pgErr.Code,
				pgErr.SQLState(),
			))
		}
		return err
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
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RecipesRepository) Update(ctx context.Context, recipe *models.Recipe) error {
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

func (r *RecipesRepository) Delete(ctx context.Context, id int) error {
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

func (r *RecipesRepository) Restore(ctx context.Context, id int) error {
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
