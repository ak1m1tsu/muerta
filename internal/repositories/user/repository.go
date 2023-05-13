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
	FindRoles(ctx context.Context, id int) ([]models.Role, error)
	FindSettings(ctx context.Context, id int) ([]models.Setting, error)
	UpdateSetting(ctx context.Context, id int, setting models.Setting) (models.Setting, error)
	DeleteStorage(ctx context.Context, id int, entity models.Storage) error
	CreateStorage(ctx context.Context, id int, entity models.Storage) (models.Storage, error)
	FindStorages(ctx context.Context, id int) ([]models.Storage, error)
}

type userRepository struct {
	client repositories.PostgresClient
}

// CreateStorage implements UserRepositorer
func (r *userRepository) CreateStorage(ctx context.Context, id int, entity models.Storage) (models.Storage, error) {
	var (
		query = `
			INSERT INTO users_storages (id_user, id_storage)
			VALUES ($1, $2)
		`
		querySelect = `
			SELECT s.id, s.name, st.name, s.temperature, s.humidity
			FROM storages s
			JOIN storages_types st ON s.id_type = st.id
			JOIN users_storages us ON us.id_storage = s.id
			WHERE us.id_user = $1
			LIMIT 1
		`
	)
	if _, err := r.client.Exec(ctx, query, id, entity.ID); err != nil {
		return models.Storage{}, fmt.Errorf("failed to create storage: %w", err)
	}
	if err := r.client.QueryRow(ctx, querySelect, id).Scan(&entity.ID, &entity.Name, &entity.Type.Name, &entity.Temperature, &entity.Humidity); err != nil {
		return models.Storage{}, fmt.Errorf("failed to scan storage: %w", err)
	}
	return entity, nil
}

// DeleteStorage implements UserRepositorer
func (r *userRepository) DeleteStorage(ctx context.Context, id int, entity models.Storage) error {
	var (
		query = `
			DELETE FROM users_storages
			WHERE id_user = $1 AND id_storage = $2
		`
	)
	if _, err := r.client.Exec(ctx, query, id, entity.ID); err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	return nil
}

// FindStorages implements UserRepositorer
func (r *userRepository) FindStorages(ctx context.Context, id int) ([]models.Storage, error) {
	var (
		query = `
			SELECT s.id, s.name, st.name, s.temperature, s.humidity
			FROM storages s
			JOIN storages_types st ON s.id_type = st.id
			JOIN users_storages us ON us.id_storage = s.id
			WHERE us.id_user = $1
		`
		entities = make([]models.Storage, 0)
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query storages: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var entity models.Storage
		if err := rows.Scan(&entity.ID, &entity.Name, &entity.Type.Name, &entity.Temperature, &entity.Humidity); err != nil {
			return nil, fmt.Errorf("failed to scan storage: %w", err)
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

// FindRoles implements UserRepositorer
func (r *userRepository) FindRoles(ctx context.Context, id int) ([]models.Role, error) {
	var (
		query = `
			SELECT r.id, r.name
			FROM roles r
			JOIN users_roles ur ON ur.id_role = r.id
			WHERE ur.id_user = $1
		`
		roles = make([]models.Role, 0)
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// FindSettings implements UserRepositorer
func (r *userRepository) FindSettings(ctx context.Context, id int) ([]models.Setting, error) {
	var (
		query = `
			SELECT s.id, s.name, us.value, sc.name
			FROM settings s
			JOIN users_settings us ON s.id = us.id_setting
			JOIN settings_categories sc ON s.id_category = sc.id
			WHERE us.id_user = $1
		`
		settings = make([]models.Setting, 0)
	)
	rows, err := r.client.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var setting models.Setting
		if err := rows.Scan(&setting.ID, &setting.Name, &setting.Value, &setting.Category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

// UpdateSetting implements UserRepositorer
func (r *userRepository) UpdateSetting(ctx context.Context, id int, setting models.Setting) (models.Setting, error) {
	var (
		query = `
			UPDATE users_settings
			SET value = $2
			WHERE id_user = $1 AND id_setting = $3
		`
		querySettings = `
			SELECT s.name, us.value, sc.name FROM settings s
			JOIN users_settings us ON s.id = us.id_setting
			JOIN settings_categories sc ON s.id_category = sc.id
			WHERE us.id_user = $1
			LIMIT 1
		`
	)
	if _, err := r.client.Exec(ctx, query, id, setting.Value, setting.ID); err != nil {
		return models.Setting{}, fmt.Errorf("failed to update setting: %w", err)
	}
	rows, err := r.client.Query(ctx, querySettings, id)
	if err != nil {
		return models.Setting{}, fmt.Errorf("failed to query settings: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&setting.Name, &setting.Value, &setting.Category.Name); err != nil {
			return models.Setting{}, fmt.Errorf("failed to scan setting: %w", err)
		}
	}
	return setting, nil
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
			JOIN users_roles ur ON ur.id_role = r.id
			JOIN users u ON u.id = ur.id_user
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
