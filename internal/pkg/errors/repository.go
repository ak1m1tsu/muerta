package errors

var (
	ErrFailedToBeginTransaction  = New("failed to begin transaction")
	ErrFailedToCommitTransaction = New("failed to commit transaction")
	ErrFailedToCoopyModels       = New("failed to copy models")
	ErrFailedToCountModels       = New("failed to count models")
)

var (
	ErrFailedToInsertUser  = New("failed to insert user")
	ErrFailedToUpdateUser  = New("failed to update user")
	ErrFailedToDeleteUser  = New("failed to delete user")
	ErrFailedToRestoreUser = New("failed to restore user")
	ErrFailedToSelectUsers = New("failed to select users")
	ErrFailedToSelectUser  = New("failed to select user")
)

var (
	ErrFailedToCreatePassword = New("failed to insert password")
	ErrFailedToSelectPassword = New("failed to select password")
)

var (
	ErrFailedToSelectRoles = New("failed to select roles")
	ErrFailedToSelectRole  = New("failed to select role")
	ErrFailedToInsertRole  = New("failed to insert role")
	ErrFailedToUpdateRole  = New("failed to update role")
	ErrFailedToDeleteRole  = New("failed to delete role")
	ErrFailedToRestoreRole = New("failed to restore role")
)

var (
	ErrFailedToSelectShelfLives = New("failed to select shelf lives")
	ErrFailedToSelectShelfLife  = New("failed to select shelf life")
	ErrFailedToInsertShelfLife  = New("failed to insert shelf life")
	ErrFailedToUpdateShelfLife  = New("failed to update shelf life")
	ErrFailedToDeleteShelfLife  = New("failed to delete shelf life")
	ErrFailedToRestoreShelfLife = New("failed to restore shelf life")
)
