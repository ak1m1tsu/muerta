package errors

import "github.com/romankravchuk/muerta/internal/pkg/errors"

var (
	ErrFailedToBeginTransaction  = errors.New("failed to begin transaction")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction")
	ErrFailedToCoopyModels       = errors.New("failed to copy models")
	ErrFailedToCountModels       = errors.New("failed to count models")
)

var (
	ErrFailedToInsertUser  = errors.New("failed to insert user")
	ErrFailedToInsertUsers = errors.New("failed to insert users")
	ErrFailedToUpdateUser  = errors.New("failed to update user")
	ErrFailedToDeleteUser  = errors.New("failed to delete user")
	ErrFailedToRestoreUser = errors.New("failed to restore user")
	ErrFailedToSelectUsers = errors.New("failed to select users")
	ErrFailedToSelectUser  = errors.New("failed to select user")
)

var (
	ErrFailedToCreatePassword = errors.New("failed to insert password")
	ErrFailedToSelectPassword = errors.New("failed to select password")
)

var (
	ErrFailedToSelectRoles = errors.New("failed to select roles")
	ErrFailedToSelectRole  = errors.New("failed to select role")
	ErrFailedToInsertRole  = errors.New("failed to insert role")
	ErrFailedToUpdateRole  = errors.New("failed to update role")
	ErrFailedToDeleteRole  = errors.New("failed to delete role")
	ErrFailedToRestoreRole = errors.New("failed to restore role")
)

var (
	ErrFailedToSelectShelfLives = errors.New("failed to select shelf lives")
	ErrFailedToSelectShelfLife  = errors.New("failed to select shelf life")
	ErrFailedToInsertShelfLife  = errors.New("failed to insert shelf life")
	ErrFailedToUpdateShelfLife  = errors.New("failed to update shelf life")
	ErrFailedToDeleteShelfLife  = errors.New("failed to delete shelf life")
	ErrFailedToRestoreShelfLife = errors.New("failed to restore shelf life")
)
