package domain

import "github.com/jhonsferg/corelog/internal/platform/apperror"

var (
	ErrInvalidName     = apperror.New(apperror.KindValidation, "name is required")
	ErrInvalidEmail    = apperror.New(apperror.KindValidation, "a valid email is required")
	ErrInvalidPassword = apperror.New(apperror.KindValidation, "password is required")
	ErrInvalidRole     = apperror.New(apperror.KindValidation, "role is invalid")

	ErrUserNotFound       = apperror.New(apperror.KindNotFound, "user not found")
	ErrEmailAlreadyTaken  = apperror.New(apperror.KindConflict, "email is already registered")
	ErrInvalidCredentials = apperror.New(apperror.KindUnauthorized, "invalid email or password")
	ErrIncorrectPassword  = apperror.New(apperror.KindUnauthorized, "current password is incorrect")
)
