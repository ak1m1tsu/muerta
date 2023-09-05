package postgres

import (
	"context"
	"database/sql"
	errs "errors"

	"github.com/romankravchuk/muerta/internal/v2/data"
	"github.com/romankravchuk/muerta/internal/v2/lib/errors"
	"github.com/romankravchuk/muerta/internal/v2/storage"
	"github.com/romankravchuk/muerta/internal/v2/storage/users"
)

// Storage represents a postgres users storage.
type Storage struct {
	db *sql.DB
}

// New creates a new postgres users storage instance.
//
// If the db pool is nil, an error storage.ErrDBPoolIsNnil is returned.
func New(db *sql.DB) (*Storage, error) {
	const op = "storage.users.postgres.New"

	if db == nil {
		return nil, errors.WithOp(op, storage.ErrDBPoolIsNnil)
	}

	return &Storage{db: db}, nil
}

// Create creates a new user in the database.
//
// The user's ID and created_on fields are populated automatically.
func (s *Storage) Create(ctx context.Context, user *data.User) error {
	const (
		op    = "storage.users.postgres.Storage.Create"
		query = `
		INSERT INTO users
			(first_name, last_name, email, encrypted_password)
		VALUES
			($1, $2, $3, $4)
		RETURNING id, created_on`
		rolesQuery = `
		INSERT INTO user_roles
			(user_id, role_id)
		SELECT $1, id
		FROM roles
		WHERE name = 'user'`
	)

	tx, err := s.db.Begin()
	if err != nil {
		return errors.WithOp(op, err)
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRowContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.EncryptedPassword,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if errs.Is(err, sql.ErrNoRows) {
			return errors.WithOp(op, users.ErrUserNotFound)
		}
		return errors.WithOp(op, err)
	}

	_, err = tx.ExecContext(ctx, rolesQuery, user.ID)
	if err != nil {
		return errors.WithOp(op, err)
	}

	if err := tx.Commit(); err != nil {
		return errors.WithOp(op, err)
	}

	return nil
}

// Update updates the user in the database.
func (s *Storage) Update(ctx context.Context, user *data.User) error {
	const (
		op    = "storage.users.postgres.Storage.Update"
		query = `
		UPDATE users
		SET first_name = $1,
			last_name = $2,
			updated_at = NOW()
		WHERE id = $3`
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.WithOp(op, err)
	}

	if _, err := stmt.ExecContext(ctx, user.FirstName, user.LastName, user.ID); err != nil {
		return errors.WithOp(op, err)
	}

	return nil
}

// Delete deletes the user from the database.
//
// Actually, the user's deleted_at field is set to now instead of actually deleting it.
func (s *Storage) Delete(ctx context.Context, id string) error {
	const (
		op    = "storage.users.postgres.Storage.Delete"
		query = `
		UPDATE users
		SET deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1`
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.WithOp(op, err)
	}

	if _, err := stmt.ExecContext(ctx, id); err != nil {
		return errors.WithOp(op, err)
	}

	return nil
}

// FindByEmail finds a user be email.
//
// If user's deleted_at field is not null, an error is returned.
func (s *Storage) FindByEmail(ctx context.Context, email string) (*data.User, error) {
	const (
		op    = "storage.users.postgres.Storage.FindByEmail"
		query = `
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.encrypted_password,
			u.created_at,
			u.updated_at
		FROM users u
		WHERE u.email = $1
		AND u.deleted_at IS NULL`
	)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.WithOp(op, err)
	}

	var user *data.User
	err = stmt.QueryRowContext(ctx, email).Scan(
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.EncryptedPassword,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		if errs.Is(err, sql.ErrNoRows) {
			return nil, errors.WithOp(op, users.ErrUserNotFound)
		}
		return nil, errors.WithOp(op, err)
	}

	return user, nil
}

// FindMany finds users by filter.
//
// For FirstName and LastName fields the filter is case-insensitive.
func (s *Storage) FindMany(ctx context.Context, filter data.UserFilter) ([]data.User, error) {
	const (
		op    = "storage.users.postgres.Storage.FindMany"
		query = `
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.encrypted_password,
			u.created_at,
			u.updated_at
		FROM users u
		WHERE u.deleted_at IS NULL
			AND u.first_name ILIKE $3
			AND u.last_name ILIKE $4
		ORDER BY u.created_at DESC
		LIMIT $1
		OFFSET $2`
	)

	foundUsers := make([]data.User, 0, filter.Limit)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.WithOp(op, err)
	}

	rows, err := stmt.QueryContext(ctx,
		filter.Limit,
		filter.Offset,
		filter.FirstName,
		filter.LastName,
	)
	if err != nil {
		return nil, errors.WithOp(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		user := data.User{}
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.EncryptedPassword,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, errors.WithOp(op, err)
		}
		foundUsers = append(foundUsers, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithOp(op, err)
	}

	if len(foundUsers) == 0 {
		return nil, errors.WithOp(op, users.ErrUsersNotFound)
	}

	return foundUsers, nil
}
