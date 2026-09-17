package domain

import "github.com/jhonsferg/corelog/internal/platform/apperror"

var (
	ErrInvalidName   = apperror.New(apperror.KindValidation, "name is required")
	ErrTeamNotFound  = apperror.New(apperror.KindNotFound, "team not found")
	ErrTeamNameTaken = apperror.New(apperror.KindConflict, "a team with this name already exists")
)
