package user

import "context"

type UserRepositorer interface {
	FindByID(ctx context.Context, id int) (any, error)
	FindByName(ctx context.Context, name string) (any, error)
	FindMany(ctx context.Context, filter any) ([]any, error)
	Create(ctx context.Context, user any) (any, error)
	Update(ctx context.Context, id int, new any) (any, error)
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
}

func New() UserRepositorer {
	return &userRepository{}
}

func (repo *userRepository) FindByID(ctx context.Context, id int) (any, error) {
	return id, nil
}

func (repo *userRepository) FindByName(ctx context.Context, name string) (any, error) {
	return name, nil
}

func (repo *userRepository) FindMany(ctx context.Context, filter any) ([]any, error) {
	return []any{filter}, nil
}

func (repo *userRepository) Create(ctx context.Context, user any) (any, error) {
	return user, nil
}

func (repo *userRepository) Update(ctx context.Context, id int, new any) (any, error) {
	return []any{id, new}, nil
}

func (repo *userRepository) Delete(ctx context.Context, id int) error {
	return nil
}
