package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

type Repository interface {
	Save(ctx context.Context, t *domain.Ticket) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error)
	FindAll(ctx context.Context, filter ListFilter) ([]*domain.Ticket, error)
	Update(ctx context.Context, t *domain.Ticket) error
}

type UserDirectory interface {
	TeamOf(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error)
}
