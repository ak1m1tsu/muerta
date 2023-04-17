package user

import (
	"bytes"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/romankravchuk/muerta/internal/repositories/models"
)

type UserRepositorer interface {
	FindByID(ctx context.Context, id int) (models.User, error)
	FindByName(ctx context.Context, name string) (models.User, error)
	FindMany(ctx context.Context, filter models.UserFilter) ([]models.User, error)
	Create(ctx context.Context, user any) (models.User, error)
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
		return models.User{}, err
	}
	return user, nil
}

func (repo *userRepository) FindByName(ctx context.Context, name string) (models.User, error) {
	var user models.User
	if err := repo.db.Get(&user, "SELECT * FROM users WHERE users.name = $1", name); err != nil {
		return models.User{}, err
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
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) Create(ctx context.Context, userDTO any) (models.User, error) {
	var (
		user   models.User
		uquery string = "INSERT INTO users (name, salt) VALUES (:name, :salt) RETURN ID"
		pquery string = "INSERT INTO passwords (passhash) VALUES (:passhash)"
		squery string = "INSERT INTO users_settings (id_user, id_setting, value) VALUES (:id_user, :id_setting, :value)"
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return models.User{}, err
	}

	_, err = tx.Exec(uquery, "roman", "test")
	if err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	_, err = tx.Exec(pquery, "test")
	if err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	_, err = tx.Exec(squery)
	if err != nil {
		tx.Rollback()
		return models.User{}, err
	}

	if err = tx.Commit(); err != nil {
		return models.User{}, err
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
