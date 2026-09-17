package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

type CreateInput struct {
	Title       string
	Description string
	Priority    domain.Priority
	RequesterID uuid.UUID
}

type ListFilter struct {
	Status   domain.Status
	Priority domain.Priority
	Page     int
	PageSize int
}

type Service interface {
	Create(ctx context.Context, in CreateInput) (*domain.Ticket, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Ticket, error)
	List(ctx context.Context, filter ListFilter) ([]*domain.Ticket, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.Status) (*domain.Ticket, error)
	Assign(ctx context.Context, id uuid.UUID, assignerID uuid.UUID, assigneeID uuid.UUID) (*domain.Ticket, error)
}
