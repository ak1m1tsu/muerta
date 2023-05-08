package user

import (
	"context"

	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/pkg/translate"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
)

type UserServicer interface {
	FindUserByID(ctx context.Context, id int) (dto.FindUserDTO, error)
	FindUsers(ctx context.Context, filter *dto.UserFilterDTO) ([]dto.FindUserDTO, error)
	CreateUser(ctx context.Context, payload *dto.CreateUserDTO) error
	UpdateUser(ctx context.Context, id int, user *dto.UpdateUserDTO) error
	DeleteUser(ctx context.Context, id int) error
	RestoreUser(ctx context.Context, id int) error
}

type userService struct {
	repo repo.UserRepositorer
}

func New(repo repo.UserRepositorer) UserServicer {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) FindUserByID(ctx context.Context, id int) (dto.FindUserDTO, error) {
	user, err := svc.repo.FindByID(ctx, id)
	result := translate.UserModelToFindDTO(&user)
	if err != nil {
		return dto.FindUserDTO{}, err
	}
	return result, nil
}

func (svc *userService) FindUsers(ctx context.Context, filter *dto.UserFilterDTO) ([]dto.FindUserDTO, error) {
	users, err := svc.repo.FindMany(ctx, filter.Limit, filter.Offset, filter.Name)
	dtos := translate.UserModelsToFindDTOs(users)
	if err != nil {
		return nil, err
	}
	return dtos, nil
}

func (svc *userService) CreateUser(ctx context.Context, payload *dto.CreateUserDTO) error {
	model := translate.CreateUserDTOToModel(payload)
	if err := svc.repo.Create(ctx, model); err != nil {
		return err
	}
	return nil
}

func (svc *userService) UpdateUser(ctx context.Context, id int, user *dto.UpdateUserDTO) error {
	oldUser, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user.Name != "" {
		oldUser.Name = user.Name
	}
	if err := svc.repo.Update(ctx, oldUser); err != nil {
		return err
	}
	return nil
}

func (svc *userService) DeleteUser(ctx context.Context, id int) error {
	if err := svc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (svc *userService) RestoreUser(ctx context.Context, id int) error {
	if err := svc.repo.Restore(ctx, id); err != nil {
		return err
	}
	return nil
}
