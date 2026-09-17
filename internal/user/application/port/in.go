package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/user/domain"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
	Role     domain.Role
}

type AuthenticateInput struct {
	Email    string
	Password string
}

type UpdateProfileInput struct {
	Name  string
	Email string
}

type ChangePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}

type SearchFilter struct {
	Query  string
	TeamID *uuid.UUID
}

type Service interface {
	Register(ctx context.Context, in RegisterInput) (*domain.User, error)
	Authenticate(ctx context.Context, in AuthenticateInput) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, in UpdateProfileInput) (*domain.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, in ChangePasswordInput) error
	SetTeam(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID) (*domain.User, error)
	SearchTeammates(ctx context.Context, actorID uuid.UUID, query string) ([]*domain.User, error)
	Search(ctx context.Context, filter SearchFilter) ([]*domain.User, error)
}
