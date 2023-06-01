package user

import (
	"context"

	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
)

type UserStorage interface {
	FindByID(ctx context.Context, id int) (models.User, error)
	FindByName(ctx context.Context, name string) (models.User, error)
	FindMany(ctx context.Context, filter models.UserFilter) ([]models.User, error)
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
	Count(ctx context.Context, filter models.UserFilter) (int, error)
	UserPasswordStorage
	UserRoleStorage
	UserVaultStorage
	UserShelfLifeStorage
	UserSettingStorage
}

type UserPasswordStorage interface {
	FindPassword(ctx context.Context, passhash string) (models.Password, error)
}

type UserRoleStorage interface {
	FindRoles(ctx context.Context, id int) ([]models.Role, error)
}

type UserVaultStorage interface {
	RemoveVault(ctx context.Context, id, storageID int) error
	AddVault(ctx context.Context, id, storageID int) (models.Vault, error)
	FindVaults(ctx context.Context, id int) ([]models.Vault, error)
}

type UserSettingStorage interface {
	FindSettings(ctx context.Context, id int) ([]models.Setting, error)
	UpdateSetting(ctx context.Context, id int, setting models.Setting) (models.Setting, error)
}

type UserShelfLifeStorage interface {
	CreateShelfLife(ctx context.Context, userId int, model models.ShelfLife) (models.ShelfLife, error)
	FindShelfLife(ctx context.Context, userId int, shelfLifeId int) (models.ShelfLife, error)
	FindShelfLives(ctx context.Context, userId int) ([]models.ShelfLife, error)
	UpdateShelfLife(ctx context.Context, userId int, model models.ShelfLife) (models.ShelfLife, error)
	DeleteShelfLife(ctx context.Context, userId int, shelfLifeId int) error
	RestoreShelfLife(ctx context.Context, userId int, shelfLifeId int) (models.ShelfLife, error)
}
