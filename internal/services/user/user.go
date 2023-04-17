package user

import (
	"context"

	"github.com/romankravchuk/muerta/internal/repositories/models"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
)

type UserServicer interface {
	FindByID(ctx context.Context, id int) (any, error)
	FindByName(ctx context.Context, name string) (any, error)
	FindMany(ctx context.Context, filter any) ([]any, error)
	Create(ctx context.Context, user any) (any, error)
	Update(ctx context.Context, id int, new any) (any, error)
	Delete(ctx context.Context, id int) error
}

type userService struct {
	repo repo.UserRepositorer
}

func New(repo repo.UserRepositorer) UserServicer {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) FindByID(ctx context.Context, id int) (any, error) {
	user, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *userService) FindByName(ctx context.Context, name string) (any, error) {
	user, err := svc.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *userService) FindMany(ctx context.Context, filter any) ([]any, error) {
	_, err := svc.repo.FindMany(ctx, models.UserFilter{})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (svc *userService) Create(ctx context.Context, userDTO any) (any, error) {
	user, err := svc.repo.Create(ctx, userDTO)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *userService) Update(ctx context.Context, id int, user any) (any, error) {
	updatedUser, err := svc.repo.Update(ctx, id, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (svc *userService) Delete(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
