package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/romankravchuk/muerta/internal/pkg/errors"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

type userStorage struct {
	c postgres.Client
}

func New(c postgres.Client) UserStorage {
	return &userStorage{c: c}
}

func (s *userStorage) Count(ctx context.Context, filter models.UserFilter) (int, error) {
	count := 0
	if err := s.c.QueryRow(ctx, countUsers, "%"+filter.Name+"%").
		Scan(&count); err != nil {
		return 0, errors.ErrFailedToCountModels.With(err)
	}
	return count, nil
}

// CreateShelfLife implements UserRepositorer
func (s *userStorage) CreateShelfLife(ctx context.Context, id int, params models.ShelfLife,
) (models.ShelfLife, error) {
	err := s.c.QueryRow(ctx, createShelfLife, id, params.Product.ID, params.Storage.ID, params.Measure.ID, params.Quantity, params.PurchaseDate, params.EndDate).
		Scan(&params.ID, &params.Product.Name, &params.Storage.Name, &params.Measure.Name)
	if err != nil {
		return models.ShelfLife{}, errors.ErrFailedToInsertShelfLife.With(err)
	}
	return params, nil
}

// DeleteShelfLife implements UserRepositorer
func (s *userStorage) DeleteShelfLife(ctx context.Context, id int, shelfLifeId int) error {
	if _, err := s.c.Exec(ctx, deleteShelfLife, id, shelfLifeId); err != nil {
		return errors.ErrFailedToDeleteShelfLife.With(err)
	}
	return nil
}

// FindShelfLife implements UserRepositorer
func (s *userStorage) FindShelfLife(ctx context.Context, id int, shelfLifeId int) (models.ShelfLife, error) {
	shelfLife := models.ShelfLife{}
	if err := s.c.QueryRow(ctx, findShelfLife, id, shelfLifeId).Scan(
		&shelfLife.ID, &shelfLife.Product.ID, &shelfLife.Storage.ID,
		&shelfLife.Measure.ID, &shelfLife.Quantity, &shelfLife.PurchaseDate,
		&shelfLife.EndDate, &shelfLife.Product.Name, &shelfLife.Storage.Name,
		&shelfLife.Measure.Name,
	); err != nil {
		return models.ShelfLife{}, errors.ErrFailedToSelectShelfLife.With(err)
	}
	return shelfLife, nil
}

// FindShelfLives implements UserRepositorer
func (s *userStorage) FindShelfLives(ctx context.Context, id int) ([]models.ShelfLife, error) {
	shelfLives := make([]models.ShelfLife, 0)
	rows, err := s.c.Query(ctx, findShelfLives, id)
	if err != nil {
		return nil, errors.ErrFailedToSelectShelfLives.With(err)
	}
	defer rows.Close()
	for rows.Next() {
		shelfLife := models.ShelfLife{}
		if err := rows.Scan(
			&shelfLife.ID, &shelfLife.Product.ID, &shelfLife.Storage.ID, &shelfLife.Measure.ID,
			&shelfLife.Quantity, &shelfLife.PurchaseDate,
			&shelfLife.EndDate, &shelfLife.Product.Name, &shelfLife.Storage.Name,
			&shelfLife.Measure.Name,
		); err != nil {
			return nil, fmt.Errorf("failed to scan shelf life: %w", err)
		}
		shelfLives = append(shelfLives, shelfLife)
	}
	return shelfLives, nil
}

// RestoreShelfLife implements UserRepositorer
func (s *userStorage) RestoreShelfLife(ctx context.Context, id, shelfLifeId int) (models.ShelfLife, error) {
	shelfLife := models.ShelfLife{}
	if err := s.c.QueryRow(ctx, restoreShelfLife, id, shelfLifeId).
		Scan(
			&shelfLife.ID,
			&shelfLife.Product.ID,
			&shelfLife.Storage.ID,
			&shelfLife.Measure.ID,
			&shelfLife.Quantity,
			&shelfLife.PurchaseDate,
			&shelfLife.EndDate,
			&shelfLife.Product.Name,
			&shelfLife.Storage.Name,
			&shelfLife.Measure.Name,
		); err != nil {
		return models.ShelfLife{}, fmt.Errorf("failed to scan shelf life: %w", err)
	}
	return shelfLife, nil
}

// UpdateShelfLife implements UserRepositorer
func (s *userStorage) UpdateShelfLife(ctx context.Context, id int, params models.ShelfLife,
) (models.ShelfLife, error) {
	if err := s.c.QueryRow(ctx, updateShelfLife, id, params.ID, params.Product.ID, params.Storage.ID, params.Measure.ID, params.Quantity, params.PurchaseDate, params.EndDate).
		Scan(
			&params.Product.ID,
			&params.Storage.ID,
			&params.Measure.ID,
			&params.Quantity,
			&params.PurchaseDate,
			&params.EndDate,
			&params.Product.Name,
			&params.Storage.Name,
			&params.Measure.Name,
		); err != nil {
		return models.ShelfLife{}, fmt.Errorf("failed to scan shelf life: %w", err)
	}
	return params, nil
}

// AddVault implements UserRepositorer
func (s *userStorage) AddVault(ctx context.Context, id, vaultId int) (models.Vault, error) {
	var vault models.Vault
	if err := s.c.QueryRow(ctx, addVault, id, vaultId).Scan(
		&vault.ID,
		&vault.Name,
		&vault.Temperature,
		&vault.Humidity,
		&vault.Type.ID,
		&vault.Type.Name,
	); err != nil {
		return models.Vault{}, fmt.Errorf("failed to scan storage: %w", err)
	}
	return vault, nil
}

// RemoveVault implements UserRepositorer
func (s *userStorage) RemoveVault(ctx context.Context, id, vaultId int) error {
	if _, err := s.c.Exec(ctx, removeVault, id, vaultId); err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	return nil
}

// FindVaults implements UserRepositorer
func (s *userStorage) FindVaults(ctx context.Context, id int) ([]models.Vault, error) {
	vaults := make([]models.Vault, 0)
	rows, err := s.c.Query(ctx, findVaults, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query storages: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		vault := models.Vault{}
		if err := rows.Scan(&vault.ID, &vault.Name, &vault.Type.Name, &vault.Temperature, &vault.Humidity); err != nil {
			return nil, fmt.Errorf("failed to scan storage: %w", err)
		}
		vaults = append(vaults, vault)
	}
	return vaults, nil
}

// FindRoles implements UserRepositorer
func (s *userStorage) FindRoles(ctx context.Context, id int) ([]models.Role, error) {
	roles := make([]models.Role, 0)
	rows, err := s.c.Query(ctx, findRoles, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		role := models.Role{}
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// FindSettings implements UserRepositorer
func (s *userStorage) FindSettings(ctx context.Context, id int) ([]models.Setting, error) {
	settings := make([]models.Setting, 0)
	rows, err := s.c.Query(ctx, findSettings, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		setting := models.Setting{}
		if err := rows.Scan(&setting.ID, &setting.Name, &setting.Value, &setting.Category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

// UpdateSetting implements UserRepositorer
func (s *userStorage) UpdateSetting(ctx context.Context, id int, setting models.Setting,
) (models.Setting, error) {
	if _, err := s.c.Exec(ctx, updateSetting, id, setting.Value, setting.ID); err != nil {
		return models.Setting{}, fmt.Errorf("failed to update setting: %w", err)
	}
	rows, err := s.c.Query(ctx, findSetting, id)
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

func (repo *userStorage) FindByID(ctx context.Context, id int) (models.User, error) {
	var (
		findUserById = `
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
	if err := repo.c.QueryRow(ctx, findUserById, id).Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
		return models.User{}, fmt.Errorf("failed to query user: %w", err)
	}
	return user, nil
}

func (repo *userStorage) FindByName(ctx context.Context, name string) (models.User, error) {
	user := models.User{Roles: make([]models.Role, 0)}
	if err := repo.c.QueryRow(ctx, findUserByName, name).Scan(&user.ID, &user.Name, &user.Salt, &user.CreatedAt); err != nil {
		return models.User{}, errors.ErrFailedToSelectUser.With(err)
	}
	rows, err := repo.c.Query(ctx, findRoles, user.ID)
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

func (repo *userStorage) FindMany(ctx context.Context, filter models.UserFilter) ([]models.User, error) {
	users := make([]models.User, 0, filter.Limit)
	rows, err := repo.c.Query(ctx, findUsers, "%"+filter.Name+"%", filter.Limit, filter.Offset)
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

func (repo *userStorage) FindPassword(ctx context.Context, hash string) (models.Password, error) {
	passwd := models.Password{}
	if err := repo.c.QueryRow(ctx, findPassword, hash).Scan(&passwd.Hash); err != nil {
		return models.Password{}, errors.ErrFailedToSelectPassword.With(err)
	}
	return passwd, nil
}

func (repo *userStorage) Create(ctx context.Context, params models.User) error {
	tx, err := repo.c.Begin(ctx)
	if err != nil {
		return errors.ErrFailedToBeginTransaction.With(err)
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, createUser, params.Name, params.Salt).
		Scan(&params.ID); err != nil {
		return errors.ErrFailedToInsertUser.With(err)
	}
	if _, err := tx.Exec(ctx, createPassword, params.Password.Hash); err != nil {
		return errors.ErrFailedToCreatePassword.With(err)
	}
	if _, err := tx.CopyFrom(ctx,
		pgx.Identifier{"users_settings"},
		[]string{"id_user", "id_setting", "value"},
		pgx.CopyFromSlice(len(params.Settings), func(i int) ([]any, error) {
			return []any{params.ID, params.Settings[i].ID, params.Settings[i].Value}, nil
		}),
	); err != nil {
		return errors.ErrFailedToCoopyModels.With(err)
	}
	if _, err := tx.CopyFrom(ctx,
		pgx.Identifier{"users_roles"},
		[]string{"id_user", "id_role"},
		pgx.CopyFromSlice(len(params.Roles), func(i int) ([]any, error) {
			return []any{params.ID, params.Roles[i].ID}, nil
		}),
	); err != nil {
		return errors.ErrFailedToCoopyModels.With(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return errors.ErrFailedToCommitTransaction.With(err)
	}
	return nil
}

func (repo *userStorage) Update(ctx context.Context, params models.User) error {
	if _, err := repo.c.Exec(ctx, updateUser, params.Name, params.ID); err != nil {
		return errors.ErrFailedToUpdateUser.With(err)
	}
	return nil
}

func (repo *userStorage) Delete(ctx context.Context, id int) error {
	if _, err := repo.c.Exec(ctx, deleteUser, id); err != nil {
		return errors.ErrFailedToDeleteUser.With(err)
	}
	return nil
}

func (repo *userStorage) Restore(ctx context.Context, id int) error {
	if _, err := repo.c.Exec(ctx, restoreUser, id); err != nil {
		return errors.ErrFailedToRestoreUser.With(err)
	}
	return nil
}
