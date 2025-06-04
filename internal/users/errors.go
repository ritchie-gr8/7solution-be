package users

import "errors"

var (
	ErrUserNotFound       = errors.New("user: not found")
	ErrEmailAlreadyExists = errors.New("user: email already exists")
	ErrInvalidCredentials = errors.New("user: invalid credentials")
	ErrUnauthorizedAccess = errors.New("user: unauthorized access")
	ErrInsertFailed       = errors.New("user: insert failed")
	ErrUpdateFailed       = errors.New("user: update failed")
	ErrDeleteFailed       = errors.New("user: delete failed")
	ErrHashingPassword    = errors.New("user: could not hash password")
	ErrGeneratingToken    = errors.New("user: could not generate token")
	ErrInvalidID          = errors.New("user: invalid ID")
)
