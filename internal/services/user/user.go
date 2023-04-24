package user

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/repositories/models"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
)

type UserServicer interface {
	FindByID(ctx context.Context, id int) (models.User, error)
	FindByName(ctx context.Context, name string) (models.User, error)
	FindMany(ctx context.Context, filter any) ([]models.User, error)
	Create(ctx context.Context, payload dto.CreateUserPayload) (models.User, error)
	Update(ctx context.Context, id int, new any) (models.User, error)
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

func (svc *userService) FindByID(ctx context.Context, id int) (models.User, error) {
	user, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (svc *userService) FindByName(ctx context.Context, name string) (models.User, error) {
	user, err := svc.repo.FindByName(ctx, name)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (svc *userService) FindMany(ctx context.Context, filter any) ([]models.User, error) {
	_, err := svc.repo.FindMany(ctx, models.UserFilter{})
	if err != nil {
		return []models.User{}, err
	}
	return nil, nil
}

func (svc *userService) Create(ctx context.Context, payload dto.CreateUserPayload) (models.User, error) {
	user := models.FilterUserPayload(payload)
	createdUser, err := svc.repo.Create(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return createdUser, nil
}

func (svc *userService) Update(ctx context.Context, id int, user any) (models.User, error) {
	updatedUser, err := svc.repo.Update(ctx, id, user)
	if err != nil {
		return models.User{}, err
	}
	return updatedUser, nil
}

func (svc *userService) Delete(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
