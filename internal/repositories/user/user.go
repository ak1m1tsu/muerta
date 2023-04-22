package user

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/romankravchuk/muerta/internal/pkg/models"
)

type UserRepositorer interface {
	FindByID(ctx context.Context, id int) (models.User, error)
	FindByName(ctx context.Context, name string) (models.User, error)
	FindMany(ctx context.Context, filter models.UserFilter) ([]models.User, error)
	FindPassword(ctx context.Context, passhash string) error
	Create(ctx context.Context, user models.User) (models.User, error)
	Update(ctx context.Context, id int, new any) (models.User, error)
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) UserRepositorer {
	return &userRepository{db: db}
}

func (repo *userRepository) FindByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	if err := repo.db.Get(&user, "SELECT * FROM users WHERE users.id = $1", id); err != nil {
		return models.User{}, fmt.Errorf("find by id: %w", err)
	}
	return user, nil
}

func (repo *userRepository) FindByName(ctx context.Context, name string) (models.User, error) {
	var user models.User
	if err := repo.db.Get(&user, "SELECT * FROM users WHERE users.name = $1", name); err != nil {
		return models.User{}, fmt.Errorf("find by name: %w", err)
	}
	return user, nil
}

func (repo *userRepository) FindMany(ctx context.Context, filter models.UserFilter) ([]models.User, error) {
	var (
		users []models.User
		query string        = "SELECT * FROM users WHERE"
		buf   *bytes.Buffer = bytes.NewBufferString(query)
	)
	if filter.Limit != 0 && filter.Offset != 0 {
		buf.WriteString(fmt.Sprintf(" OFFSET %d ROWS FETCH FIRST %d ROWS ONLY ", filter.Offset, filter.Limit))
	}
	if err := repo.db.Select(&users, buf.String()); err != nil {
		return nil, fmt.Errorf("find many: %w", err)
	}
	return users, nil
}

func (repo *userRepository) FindPassword(ctx context.Context, passhash string) error {
	fmt.Println(passhash)
	var (
		pwd   string
		query string = "SELECT * FROM passwords WHERE passwords.passhash = $1 LIMIT 1"
	)
	if err := repo.db.Get(&pwd, query, passhash); err != nil {
		return fmt.Errorf("find password: %w", err)
	}
	return nil
}

func (repo *userRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	var (
		uquery string = "INSERT INTO users (name, salt) VALUES (:name, :salt) RETURN ID"
		pquery string = "INSERT INTO passwords (passhash) VALUES (:passhash)"
		squery string = "INSERT INTO users_settings (id_user, id_setting, value) VALUES (:id_user, :id_setting, :value)"
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return models.User{}, fmt.Errorf("create: %w", err)
	}

	_, err = tx.Exec(uquery, user.Name, user.Salt)
	if err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("create: %w", err)
	}

	_, err = tx.Exec(pquery, user.Password.Hash)
	if err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("create: %w", err)
	}

	_, err = tx.Exec(squery)
	if err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("create: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return models.User{}, fmt.Errorf("create: %w", err)
	}

	return user, nil
}

func (repo *userRepository) Update(ctx context.Context, id int, new any) (models.User, error) {
	var user models.User
	return user, nil
}

func (repo *userRepository) Delete(ctx context.Context, id int) error {
	return nil
}
