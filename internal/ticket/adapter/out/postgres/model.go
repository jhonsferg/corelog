package postgres

import (
	"time"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

type ticketRow struct {
	ID          uuid.UUID  `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	Priority    string     `db:"priority"`
	RequesterID uuid.UUID  `db:"requester_id"`
	AssigneeID  *uuid.UUID `db:"assignee_id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (r ticketRow) toDomain() *domain.Ticket {
	return &domain.Ticket{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Status:      domain.Status(r.Status),
		Priority:    domain.Priority(r.Priority),
		RequesterID: r.RequesterID,
		AssigneeID:  r.AssigneeID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func fromDomain(t *domain.Ticket) ticketRow {
	return ticketRow{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		RequesterID: t.RequesterID,
		AssigneeID:  t.AssigneeID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
