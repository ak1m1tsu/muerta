package role

import (
	"context"
	"fmt"

	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/services/utils"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/role"
)

type RoleServicer interface {
	FindRoleByID(ctx context.Context, id int) (params.FindRole, error)
	FindRoles(ctx context.Context, filter *params.RoleFilter) ([]params.FindRole, error)
	CreateRole(ctx context.Context, payload *params.CreateRole) error
	UpdateRole(ctx context.Context, id int, payload *params.UpdateRole) error
	DeleteRole(ctx context.Context, id int) error
	RestoreRole(ctx context.Context, id int) error
	Count(ctx context.Context, filter params.RoleFilter) (int, error)
}

type roleService struct {
	repo repository.RoleRepositorer
}

func New(repo repository.RoleRepositorer) RoleServicer {
	return &roleService{
		repo: repo,
	}
}

func (s *roleService) Count(ctx context.Context, filter params.RoleFilter) (int, error) {
	count, err := s.repo.Count(ctx, models.RoleFilter{Name: filter.Name})
	if err != nil {
		return 0, fmt.Errorf("error counting roles: %w", err)
	}
	return count, nil
}

// CreateRole implements RoleServicer
func (s *roleService) CreateRole(ctx context.Context, payload *params.CreateRole) error {
	model := utils.CreateRoleToModel(payload)
	if err := s.repo.Create(ctx, model); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

// DeleteRole implements RoleServicer
func (s *roleService) DeleteRole(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// FindRoleByID implements RoleServicer
func (s *roleService) FindRoleByID(ctx context.Context, id int) (params.FindRole, error) {
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return params.FindRole{}, fmt.Errorf("failed to find role: %w", err)
	}
	dto := utils.RoleModelToFindRole(&role)
	return dto, nil
}

// FindRoles implements RoleServicer
func (s *roleService) FindRoles(
	ctx context.Context,
	filter *params.RoleFilter,
) ([]params.FindRole, error) {
	roles, err := s.repo.FindMany(ctx, models.RoleFilter{
		PageFilter: models.PageFilter{
			Limit:  filter.Limit,
			Offset: filter.Offset,
		},
		Name: filter.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find roles: %w", err)
	}
	dtos := utils.RoleModelsToFindRoles(roles)
	return dtos, nil
}

// RestoreRole implements RoleServicer
func (s *roleService) RestoreRole(ctx context.Context, id int) error {
	if err := s.repo.Restore(ctx, id); err != nil {
		return fmt.Errorf("failed to restore role: %w", err)
	}
	return nil
}

// UpdateRole implements RoleServicer
func (s *roleService) UpdateRole(ctx context.Context, id int, payload *params.UpdateRole) error {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}
	if payload.Name != "" {
		model.Name = payload.Name
	}
	if err := s.repo.Update(ctx, model); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}
