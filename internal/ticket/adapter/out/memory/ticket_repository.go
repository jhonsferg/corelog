package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/ticket/application/port"
	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

type TicketRepository struct {
	mu      sync.RWMutex
	tickets map[uuid.UUID]*domain.Ticket
}

func NewTicketRepository() *TicketRepository {
	repo := &TicketRepository{tickets: make(map[uuid.UUID]*domain.Ticket)}
	for _, t := range seedTickets() {
		repo.tickets[t.ID] = t
	}
	return repo
}

var _ port.Repository = (*TicketRepository)(nil)

func (r *TicketRepository) Save(_ context.Context, t *domain.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tickets[t.ID] = cloneTicket(t)
	return nil
}

func (r *TicketRepository) FindByID(_ context.Context, id uuid.UUID) (*domain.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tickets[id]
	if !ok {
		return nil, domain.ErrTicketNotFound
	}
	return cloneTicket(t), nil
}

func (r *TicketRepository) FindAll(_ context.Context, filter port.ListFilter) ([]*domain.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	matches := make([]*domain.Ticket, 0, len(r.tickets))
	for _, t := range r.tickets {
		if filter.Status != "" && t.Status != filter.Status {
			continue
		}
		if filter.Priority != "" && t.Priority != filter.Priority {
			continue
		}
		matches = append(matches, cloneTicket(t))
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].CreatedAt.After(matches[j].CreatedAt)
	})

	start := (filter.Page - 1) * filter.PageSize
	if start < 0 || start >= len(matches) {
		return []*domain.Ticket{}, nil
	}
	end := start + filter.PageSize
	if end > len(matches) {
		end = len(matches)
	}

	return matches[start:end], nil
}

func (r *TicketRepository) Update(_ context.Context, t *domain.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tickets[t.ID]; !ok {
		return domain.ErrTicketNotFound
	}
	r.tickets[t.ID] = cloneTicket(t)
	return nil
}

func cloneTicket(t *domain.Ticket) *domain.Ticket {
	clone := *t
	if t.AssigneeID != nil {
		id := *t.AssigneeID
		clone.AssigneeID = &id
	}
	return &clone
}

func seedTickets() []*domain.Ticket {
	requester := uuid.MustParse("22222222-2222-4222-8222-222222222222")

	return []*domain.Ticket{
		{
			ID:          uuid.MustParse("11111111-1111-4111-8111-111111111111"),
			Title:       "Production API returning 500s",
			Description: "The /orders endpoint intermittently returns HTTP 500 under load.",
			Status:      domain.StatusOpen,
			Priority:    domain.PriorityCritical,
			RequesterID: requester,
			CreatedAt:   time.Date(2026, 9, 10, 9, 30, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2026, 9, 10, 9, 30, 0, 0, time.UTC),
		},
		{
			ID:          uuid.MustParse("33333333-3333-4333-8333-333333333333"),
			Title:       "VPN connection drops for remote staff",
			Description: "Several employees report the VPN disconnecting every ~15 minutes.",
			Status:      domain.StatusInProgress,
			Priority:    domain.PriorityHigh,
			RequesterID: requester,
			CreatedAt:   time.Date(2026, 9, 12, 14, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2026, 9, 13, 8, 15, 0, 0, time.UTC),
		},
		{
			ID:          uuid.MustParse("44444444-4444-4444-8444-444444444444"),
			Title:       "Typo on the billing invoice template",
			Description: "\"Recieved\" should be \"Received\" on the PDF invoice header.",
			Status:      domain.StatusResolved,
			Priority:    domain.PriorityLow,
			RequesterID: requester,
			CreatedAt:   time.Date(2026, 9, 5, 11, 20, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2026, 9, 6, 10, 0, 0, 0, time.UTC),
		},
	}
}
