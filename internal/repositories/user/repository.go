package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/errors"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type UserRepositorer interface {
	FindByID(ctx context.Context, id int) (models.User, error)
	FindByName(ctx context.Context, name string) (models.User, error)
	FindMany(ctx context.Context, filter models.UserFilter) ([]models.User, error)
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
	FindShelfLives(ctx context.Context, userId int) ([]models.ShelfLife, error)
	CreateShelfLife(ctx context.Context, userId int, model models.ShelfLife) (models.ShelfLife, error)
	FindShelfLife(ctx context.Context, userId int, shelfLifeId int) (models.ShelfLife, error)
	UpdateShelfLife(ctx context.Context, userId int, model models.ShelfLife) (models.ShelfLife, error)
	DeleteShelfLife(ctx context.Context, userId int, shelfLifeId int) error
	RestoreShelfLife(ctx context.Context, userId int, shelfLifeId int) (models.ShelfLife, error)
	Count(ctx context.Context, filter models.UserFilter) (int, error)
}

type userRepository struct {
	client repositories.PostgresClient
}

func New(client repositories.PostgresClient) UserRepositorer {
	return &userRepository{client: client}
}

func (r *userRepository) Count(ctx context.Context, filter models.UserFilter) (int, error) {
	var (
		query = `
			SELECT COUNT(*) 
			FROM users 
			WHERE deleted_at IS NULL AND
				name ILIKE $1
		`
		count int
	)
	if err := r.client.QueryRow(ctx, query, "%"+filter.Name+"%").Scan(&count); err != nil {
		return 0, errors.ErrFailedToCountModels.With(err)
	}
	return count, nil
}

// CreateShelfLife implements UserRepositorer
func (r *userRepository) CreateShelfLife(ctx context.Context, userId int, model models.ShelfLife) (models.ShelfLife, error) {
	var (
		query = `
			WITH inserted AS (
				INSERT INTO shelf_lives (id_user, id_product, id_storage, id_measure, quantity, purchase_date, end_date)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id, id_product, id_storage, id_measure, quantity, purchase_date, end_date
			)
			SELECT 
				i.id, p.name, s.name, m.name
			FROM inserted i
			JOIN products p ON i.id_product = p.id
			JOIN storages s ON i.id_storage = s.id
			JOIN measures m ON i.id_measure = m.id
			WHERE p.deleted_at IS NULL AND 
				s.deleted_at IS NULL
			LIMIT 1
		`
	)
	err := r.client.QueryRow(ctx, query, userId, model.Product.ID, model.Storage.ID, model.Measure.ID, model.Quantity, model.PurchaseDate, model.EndDate).
		Scan(&model.ID, &model.Product.Name, &model.Storage.Name, &model.Measure.Name)
	if err != nil {
		return models.ShelfLife{}, errors.ErrFailedToInsertShelfLife.With(err)
	}
	return model, nil
}

// DeleteShelfLife implements UserRepositorer
func (r *userRepository) DeleteShelfLife(ctx context.Context, userId int, shelfLifeId int) error {
	var (
		query = `
			UPDATE shelf_lives
			SET deleted_at = NOW(),
				updated_at = NOW()
			WHERE id_user = $1 AND id = $2
		`
	)
	if _, err := r.client.Exec(ctx, query, userId, shelfLifeId); err != nil {
		return errors.ErrFailedToDeleteShelfLife.With(err)
	}
	return nil
}

// FindShelfLife implements UserRepositorer
func (r *userRepository) FindShelfLife(ctx context.Context, userId int, shelfLifeId int) (models.ShelfLife, error) {
	var (
		query = `
			SELECT 
				sl.id, sl.id_product, sl.id_storage, sl.id_measure, 
				sl.quantity, sl.purchase_date, sl.end_date,
				p.name, s.name, m.name
			FROM shelf_lives sl
			JOIN products p ON sl.id_product = p.id
			JOIN storages s ON sl.id_storage = s.id
			JOIN measures m ON sl.id_measure = m.id
			WHERE sl.id_user = $1 AND 
				sl.id = $2 AND
				sl.deleted_at IS NULL
			LIMIT 1
		`
		model models.ShelfLife
	)
	if err := r.client.QueryRow(ctx, query, userId, shelfLifeId).Scan(
		&model.ID, &model.Product.ID, &model.Storage.ID,
		&model.Measure.ID, &model.Quantity, &model.PurchaseDate,
		&model.EndDate, &model.Product.Name, &model.Storage.Name,
		&model.Measure.Name,
	); err != nil {
		return models.ShelfLife{}, errors.ErrFailedToSelectShelfLife.With(err)
	}
	return model, nil
}

// FindShelfLives implements UserRepositorer
func (r *userRepository) FindShelfLives(ctx context.Context, userId int) ([]models.ShelfLife, error) {
	var (
		query = `
			SELECT 
				sl.id, sl.id_product, sl.id_storage, sl.id_measure, 
				sl.quantity, sl.purchase_date, sl.end_date,
				p.name, s.name, m.name
			FROM shelf_lives sl
			JOIN products p ON sl.id_product = p.id
			JOIN storages s ON sl.id_storage = s.id
			JOIN measures m ON sl.id_measure = m.id
			WHERE sl.id_user = $1 AND 
				sl.deleted_at IS NULL
			ORDER BY sl.end_date DESC
		`
		entities []models.ShelfLife
	)
	rows, err := r.client.Query(ctx, query, userId)
	if err != nil {
		return nil, errors.ErrFailedToSelectShelfLives.With(err)
	}
	defer rows.Close()
	for rows.Next() {
		var entity models.ShelfLife
		if err := rows.Scan(
			&entity.ID, &entity.Product.ID, &entity.Storage.ID, &entity.Measure.ID,
			&entity.Quantity, &entity.PurchaseDate,
			&entity.EndDate, &entity.Product.Name, &entity.Storage.Name,
			&entity.Measure.Name,
		); err != nil {
			return nil, fmt.Errorf("failed to scan shelf life: %w", err)
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

// RestoreShelfLife implements UserRepositorer
func (r *userRepository) RestoreShelfLife(ctx context.Context, userId int, shelfLifeId int) (models.ShelfLife, error) {
	var (
		query = `
			WITH updated AS (
				UPDATE shelf_lives
				SET deleted_at = NULL,
					updated_at = NOW()
				WHERE id_user = $1 AND id = $2
				RETURNING id, id_product, id_storage, id_measure, quantity, purchase_date, end_date
			)
			SELECT 
				u.id, u.id_product, u.id_storage, u.id_measure, 
				u.quantity, u.purchase_date, u.end_date,
				p.name, s.name, m.name
			FROM updated u
			JOIN products p ON u.id_product = p.id
			JOIN storages s ON u.id_storage = s.id
			JOIN measures m ON u.id_measure = m.id
			WHERE p.deleted_at IS NULL AND
				s.deleted_at IS NULL AND
			LIMIT 1
		`
		model models.ShelfLife
	)
	if err := r.client.QueryRow(ctx, query, userId, shelfLifeId).
		Scan(
			&model.ID,
			&model.Product.ID,
			&model.Storage.ID,
			&model.Measure.ID,
			&model.Quantity,
			&model.PurchaseDate,
			&model.EndDate,
			&model.Product.Name,
			&model.Storage.Name,
			&model.Measure.Name,
		); err != nil {
		return models.ShelfLife{}, fmt.Errorf("failed to scan shelf life: %w", err)
	}
	return model, nil
}

// UpdateShelfLife implements UserRepositorer
func (r *userRepository) UpdateShelfLife(ctx context.Context, userId int, model models.ShelfLife) (models.ShelfLife, error) {
	var (
		query = `
			WITH updated AS (
				UPDATE shelf_lives
				SET id_product = $3,
					id_storage = $4,
					id_measure = $5,
					quantity = $6,
					purchase_date = $7,
					end_date = $8,
					updated_at = NOW()
				WHERE id_user = $1 AND id = $2
				RETURNING id_product, id_storage, id_measure, quantity, purchase_date, end_date
			)
			SELECT 
				u.id_product, u.id_storage, u.id_measure, 
				u.quantity, u.purchase_date, u.end_date,
				p.name, s.name, m.name
			FROM updated u
			JOIN products p ON u.id_product = p.id
			JOIN storages s ON u.id_storage = s.id
			JOIN measures m ON u.id_measure = m.id
			WHERE p.deleted_at IS NULL AND
				s.deleted_at IS NULL AND
			LIMIT 1
		`
	)
	if err := r.client.QueryRow(ctx, query, userId, model.ID, model.Product.ID, model.Storage.ID, model.Measure.ID, model.Quantity, model.PurchaseDate, model.EndDate).
		Scan(
			&model.Product.ID,
			&model.Storage.ID,
			&model.Measure.ID,
			&model.Quantity,
			&model.PurchaseDate,
			&model.EndDate,
			&model.Product.Name,
			&model.Storage.Name,
			&model.Measure.Name,
		); err != nil {
		return models.ShelfLife{}, fmt.Errorf("failed to scan shelf life: %w", err)
	}
	return model, nil
}

// CreateStorage implements UserRepositorer
func (r *userRepository) CreateStorage(ctx context.Context, id int, entity models.Storage) (models.Storage, error) {
	var (
		query = `
			WITH inserted AS (
				INSERT INTO users_storages (id_user, id_storage)
				VALUES ($1, $2)
				RETURNING id_user, id_storage
			)
			SELECT s.id, s.name, s.temperature, s.humidity, st.id, st.name
			FROM storages s
			JOIN storages_types st ON s.id_type = st.id
			JOIN inserted i ON i.id_storage = s.id
			WHERE i.id_user = $1 AND 
				s.id = i.id_storage AND
				s.deleted_at IS NULL
			LIMIT 1
		`
	)
	if err := r.client.QueryRow(ctx, query, id, entity.ID).Scan(
		&entity.ID,
		&entity.Name,
		&entity.Temperature,
		&entity.Humidity,
		&entity.Type.ID,
		&entity.Type.Name,
	); err != nil {
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

func (repo *userRepository) FindByID(ctx context.Context, id int) (models.User, error) {
	var (
		query = `
			SELECT id, name, created_at
			FROM users
			WHERE id = $1
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
		return models.User{}, errors.ErrFailedToSelectUser.With(err)
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

func (repo *userRepository) FindMany(ctx context.Context, filter models.UserFilter) ([]models.User, error) {
	var (
		query = `
			SELECT id, name, created_at
			FROM users
			WHERE name LIKE $1 AND deleted_at IS NULL
			ORDER BY created_at DESC
			LIMIT $2
			OFFSET $3
		`
		users = make([]models.User, 0, filter.Limit)
	)
	rows, err := repo.client.Query(ctx, query, "%"+filter.Name+"%", filter.Limit, filter.Offset)
	if err != nil {
		return nil, errors.ErrFailedToSelectUsers.With(err)
	}
	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
			return nil, errors.ErrFailedToSelectUser.With(err)
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
		return models.Password{}, errors.ErrFailedToSelectPassword.With(err)
	}
	return password, nil
}

func (repo *userRepository) Create(ctx context.Context, user models.User) error {
	var (
		err   error
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
	tx, err := repo.client.Begin(ctx)
	if err != nil {
		return errors.ErrFailedToBeginTransaction.With(err)
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, query, user.Name, user.Salt).Scan(&user.ID); err != nil {
		return errors.ErrFailedToInsertUser.With(err)
	}
	if _, err := tx.Exec(ctx, queryPassword, user.Password.Hash); err != nil {
		return errors.ErrFailedToCreatePassword.With(err)
	}
	if _, err := tx.CopyFrom(ctx,
		pgx.Identifier{"users_settings"},
		[]string{"id_user", "id_setting", "value"},
		pgx.CopyFromSlice(len(user.Settings), func(i int) ([]any, error) {
			return []any{user.ID, user.Settings[i].ID, user.Settings[i].Value}, nil
		}),
	); err != nil {
		return errors.ErrFailedToCoopyModels.With(err)
	}
	if _, err := tx.CopyFrom(ctx,
		pgx.Identifier{"users_roles"},
		[]string{"id_user", "id_role"},
		pgx.CopyFromSlice(len(user.Roles), func(i int) ([]any, error) {
			return []any{user.ID, user.Roles[i].ID}, nil
		}),
	); err != nil {
		return errors.ErrFailedToCoopyModels.With(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return errors.ErrFailedToCommitTransaction.With(err)
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
		return errors.ErrFailedToUpdateUser.With(err)
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
		return errors.ErrFailedToDeleteUser.With(err)
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
		return errors.ErrFailedToRestoreUser.With(err)
	}
	return nil
}
