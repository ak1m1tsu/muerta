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
	FindByName(ctx context.Context, name string) (models.Recipe, error)
	FindMany(ctx context.Context) ([]models.Recipe, error)
	Create(ctx context.Context, recipe *models.Recipe) error
	Update(ctx context.Context, recipe *models.Recipe) error
	Delete(ctx context.Context, id int) error
}

type RecipesRepository struct {
	client repositories.PostgresClient
}

func New(client repositories.PostgresClient) *RecipesRepository {
	return &RecipesRepository{client: client}
}

func (r *RecipesRepository) FindByID(ctx context.Context, id int) (models.Recipe, error) {
	return models.Recipe{}, nil
}

func (r *RecipesRepository) FindByName(ctx context.Context, name string) (models.Recipe, error) {
	return models.Recipe{}, nil
}

func (r *RecipesRepository) FindMany(ctx context.Context) ([]models.Recipe, error) {
	return []models.Recipe{}, nil
}

func (r *RecipesRepository) Create(ctx context.Context, recipe *models.Recipe) error {
	var (
		queryInsertRecipe = `
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
	if err := tx.QueryRow(ctx, queryInsertRecipe, recipe.Name, recipe.Description).Scan(&recipe.ID); err != nil {
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
			return []any{recipe.ID, recipe.Steps[i].StepID, recipe.Steps[i].Place}, nil
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
	return nil
}

func (r *RecipesRepository) Delete(ctx context.Context, id int) error {
	return nil
}
