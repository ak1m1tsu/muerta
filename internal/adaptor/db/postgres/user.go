package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/romankravchuk/muerta/internal/domain/entity"
	"github.com/romankravchuk/nix/postgres"
)

type UserStorage struct {
	db *postgres.Postgres
}

func NewUserStorage(db *postgres.Postgres) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) Create(ctx context.Context, user entity.User) (entity.User, error) {
	const sql = `
		INSERT INTO users (
			username,
			email,
			hashed_password
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING user_id, created_at
	`

	args := []any{
		user.Name,
		user.Email,
		user.Password.Hash,
	}

	err := s.db.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&user.ID,
			&user.CreatedAt,
		)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *UserStorage) GetAll(ctx context.Context, limit, offset int) (entity.Users, error) {
	const sql = `
		SELECT
			user_id,
			username,
			email,
			hashed_password,
			created_at,
			is_deleted
		FROM users
		LIMIT $1
		OFFSET $2
	`

	rows, err := s.db.Pool.Query(ctx, sql, limit, offset)
	if err != nil {
		return entity.Users{}, err
	}
	defer rows.Close()

	var users entity.Users
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password.Hash,
			&user.CreatedAt,
			&user.IsDeleted,
		); err != nil {
			return entity.Users{}, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return entity.Users{}, err
	}

	return users, nil
}

func (s *UserStorage) GetByID(ctx context.Context, id string) (entity.User, error) {
	const sql = `
		SELECT
			user_id,
			username,
			email,
			hashed_password,
			created_at,
			is_deleted
		FROM users
		WHERE user_id = $1
	`

	return s.getOne(ctx, sql, id)
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	const sql = `
		SELECT
			user_id,
			username,
			email,
			hashed_password,
			created_at,
			is_deleted
		FROM users
		WHERE email = $1
	`

	return s.getOne(ctx, sql, email)
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	const sql = `
		SELECT
			user_id,
			username,
			email,
			hashed_password,
			created_at,
			is_deleted
		FROM users
		WHERE username = $1
	`

	return s.getOne(ctx, sql, username)
}

func (s *UserStorage) getOne(ctx context.Context, sql string, args ...any) (entity.User, error) {
	var user entity.User

	err := s.db.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password.Hash,
			&user.CreatedAt,
			&user.IsDeleted,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}

func (s *UserStorage) Delete(ctx context.Context, id string) error {
	const sql = `
		UPDATE users
		SET is_deleted = true
		WHERE user_id = $1
	`

	_, err := s.db.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}
