package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/team/domain"
)

type Service interface {
	Create(ctx context.Context, name string) (*domain.Team, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Team, error)
	List(ctx context.Context) ([]*domain.Team, error)
	Rename(ctx context.Context, id uuid.UUID, name string) (*domain.Team, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
