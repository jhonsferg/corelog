package domain

import "github.com/jhonsferg/corelog/internal/platform/apperror"

var (
	ErrInvalidTitle       = apperror.New(apperror.KindValidation, "title is required")
	ErrInvalidDescription = apperror.New(apperror.KindValidation, "description is required")
	ErrInvalidPriority    = apperror.New(apperror.KindValidation, "priority is invalid")
	ErrInvalidStatus      = apperror.New(apperror.KindValidation, "status is invalid")
	ErrInvalidRequester   = apperror.New(apperror.KindValidation, "requester is required")

	ErrTicketNotFound      = apperror.New(apperror.KindNotFound, "ticket not found")
	ErrIllegalTransition   = apperror.New(apperror.KindConflict, "illegal ticket status transition")
	ErrRequesterNotFound   = apperror.New(apperror.KindValidation, "requester does not exist")
	ErrAssigneeNotFound    = apperror.New(apperror.KindValidation, "assignee does not exist")
	ErrCrossTeamAssignment = apperror.New(apperror.KindForbidden, "the assignee must belong to the same team as the assigner")
)
