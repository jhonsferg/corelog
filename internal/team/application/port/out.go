package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/team/domain"
)

type Repository interface {
	Save(ctx context.Context, t *domain.Team) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Team, error)
	FindByName(ctx context.Context, name string) (*domain.Team, error)
	FindAll(ctx context.Context) ([]*domain.Team, error)
	Update(ctx context.Context, t *domain.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
}
