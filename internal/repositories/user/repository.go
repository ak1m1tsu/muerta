package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type UserRepositorer interface {
	FindByID(ctx context.Context, id int) (models.User, error)
	FindByName(ctx context.Context, name string) (models.User, error)
	FindMany(ctx context.Context, limit, offset int, name string) ([]models.User, error)
	FindPassword(ctx context.Context, passhash string) (models.Password, error)
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	Exists(ctx context.Context, name string) bool
}

type userRepository struct {
	client repositories.PostgresClient
}

// Exists implements UserRepositorer
func (r *userRepository) Exists(ctx context.Context, name string) bool {
	if _, err := r.FindByName(ctx, name); err == nil {
		return true
	}
	return false
}

func New(client repositories.PostgresClient) UserRepositorer {
	return &userRepository{client: client}
}

func (repo *userRepository) FindByID(ctx context.Context, id int) (models.User, error) {
	var (
		query = `
			SELECT id, name, created_at
			FROM users
			WHERE id = $1
			LIMIT 1
		`
		querySettings = `
			SELECT s.name, us.value, sc.name FROM settings s
			JOIN users_settings us ON s.id = us.id_setting
			JOIN settings_categories sc ON s.id_category = sc.id
			WHERE us.id_user = $1
		`
		queryRoles = `
			SELECT r.id, r.name
			FROM roles r
			JOIN users_roles ur ON ur.id_role == r.id
			JOIN users u ON u.id == ur.id_user
			WHERE ur.id_user = $1
			LIMIT 1
		`
		user = models.User{
			Settings: make([]models.Setting, 0),
			Roles:    make([]models.Role, 0),
		}
	)
	if err := repo.client.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
		return models.User{}, fmt.Errorf("failed to query user: %w", err)
	}
	rows, err := repo.client.Query(ctx, querySettings, id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to query settings: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		setting := models.Setting{Category: models.SettingCategory{}}
		if err := rows.Scan(&setting.Name, &setting.Value, &setting.Category.Name); err != nil {
			return models.User{}, fmt.Errorf("failed to scan setting: %w", err)
		}
		user.Settings = append(user.Settings, setting)
	}
	rows, err = repo.client.Query(ctx, queryRoles, user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		role := models.Role{}
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return models.User{}, fmt.Errorf("failed to query role: %w", err)
		}
		user.Roles = append(user.Roles, role)
	}
	return user, nil
}

func (repo *userRepository) FindByName(ctx context.Context, name string) (models.User, error) {
	var (
		query = `
			SELECT id, name, salt, created_at
			FROM users
			WHERE name = $1
			LIMIT 1
		`
		queryRoles = `
			SELECT r.id, r.name
			FROM roles r
			JOIN users_roles ur ON ur.id_role = r.id
			JOIN users u ON u.id = ur.id_user
			WHERE ur.id_user = $1
			LIMIT 1
		`
		user = models.User{
			Roles: make([]models.Role, 0),
		}
	)
	if err := repo.client.QueryRow(ctx, query, name).Scan(&user.ID, &user.Name, &user.Salt, &user.CreatedAt); err != nil {
		return models.User{}, fmt.Errorf("failed to query user: %w", err)
	}
	rows, err := repo.client.Query(ctx, queryRoles, user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		role := models.Role{}
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return models.User{}, fmt.Errorf("failed to query role: %w", err)
		}
		user.Roles = append(user.Roles, role)
	}
	return user, nil
}

func (repo *userRepository) FindMany(ctx context.Context, limit, offset int, name string) ([]models.User, error) {
	var (
		query = `
			SELECT id, name, created_at
			FROM users
			WHERE name LIKE $1
			ORDER BY created_at DESC
			LIMIT $2
			OFFSET $3
		`
		users = make([]models.User, 0, limit)
	)
	rows, err := repo.client.Query(ctx, query, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		fmt.Println(user)
		users = append(users, user)
	}
	return users, nil
}

func (repo *userRepository) FindPassword(ctx context.Context, passhash string) (models.Password, error) {
	var (
		query = `
			SELECT passhash
			FROM passwords
			WHERE passhash = $1
			LIMIT 1
		`
		password models.Password
	)
	if err := repo.client.QueryRow(ctx, query, passhash).Scan(&password.Hash); err != nil {
		return models.Password{}, fmt.Errorf("failed to query password: %w", err)
	}
	return password, nil
}

func (repo *userRepository) Create(ctx context.Context, user models.User) error {
	var (
		err   error
		buf   strings.Builder
		query = `
			INSERT INTO users 
				(name, salt)
			VALUES
				($1, $2)
			RETURNING id
		`
		queryPassword = `
			INSERT INTO passwords (passhash)
			VALUES ($1)
		`
	)
	for i, s := range user.Settings {
		if i == len(user.Settings)-1 {
			_, err = buf.WriteString(fmt.Sprintf("(%d, %d, '%s')", user.ID, s.ID, s.Value))
		} else {
			_, err = buf.WriteString(fmt.Sprintf("(%d, %d, '%s'), ", user.ID, s.ID, s.Value))
		}
		if err != nil {
			return fmt.Errorf("failed to write settings: %w", err)
		}
	}
	tx, err := repo.client.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, query, user.Name, user.Salt).Scan(&user.ID); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	if _, err := tx.Exec(ctx, queryPassword, user.Password.Hash); err != nil {
		return fmt.Errorf("failed to insert password: %w", err)
	}
	if _, err := tx.CopyFrom(ctx,
		pgx.Identifier{"users_settings"},
		[]string{"id_user", "id_setting", "value"},
		pgx.CopyFromSlice(len(user.Settings), func(i int) ([]any, error) {
			return []any{user.ID, user.Settings[i].ID, user.Settings[i].Value}, nil
		}),
	); err != nil {
		return fmt.Errorf("failed to copy settings: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (repo *userRepository) Update(ctx context.Context, user models.User) error {
	var (
		query = `
			UPDATE users
			SET name = $1,
				updated_at = NOW()
			WHERE id = $2
		`
	)
	if _, err := repo.client.Exec(ctx, query, user.Name, user.ID); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (repo *userRepository) Delete(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE users
			SET deleted_at = NOW()
			SET updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := repo.client.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Restore(ctx context.Context, id int) error {
	var (
		query = `
			UPDATE users
			SET deleted_at = NULL
			SET updated_at = NOW()
			WHERE id = $1
		`
	)
	if _, err := repo.client.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}
