package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/ticket/application/port"
	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

type TicketService struct {
	repo      port.Repository
	directory port.UserDirectory
}

func NewTicketService(repo port.Repository, directory port.UserDirectory) *TicketService {
	return &TicketService{repo: repo, directory: directory}
}

var _ port.Service = (*TicketService)(nil)

func (s *TicketService) Create(ctx context.Context, in port.CreateInput) (*domain.Ticket, error) {
	ticket, err := domain.NewTicket(in.Title, in.Description, in.Priority, in.RequesterID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *TicketService) Get(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TicketService) List(ctx context.Context, filter port.ListFilter) ([]*domain.Ticket, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}
	return s.repo.FindAll(ctx, filter)
}

func (s *TicketService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.Status) (*domain.Ticket, error) {
	ticket, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := ticket.ChangeStatus(status); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *TicketService) Assign(ctx context.Context, id uuid.UUID, assignerID uuid.UUID, assigneeID uuid.UUID) (*domain.Ticket, error) {
	ticket, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	assignerTeam, err := s.directory.TeamOf(ctx, assignerID)
	if err != nil {
		return nil, err
	}

	assigneeTeam, err := s.directory.TeamOf(ctx, assigneeID)
	if err != nil {
		return nil, err
	}

	if assignerTeam == nil || assigneeTeam == nil || *assignerTeam != *assigneeTeam {
		return nil, domain.ErrCrossTeamAssignment
	}

	if err := ticket.AssignTo(assigneeID); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
