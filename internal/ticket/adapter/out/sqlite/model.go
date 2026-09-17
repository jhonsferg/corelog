package sqlite

import (
	"time"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

type ticketRow struct {
	ID          string  `db:"id"`
	Title       string  `db:"title"`
	Description string  `db:"description"`
	Status      string  `db:"status"`
	Priority    string  `db:"priority"`
	RequesterID string  `db:"requester_id"`
	AssigneeID  *string `db:"assignee_id"`
	CreatedAt   string  `db:"created_at"`
	UpdatedAt   string  `db:"updated_at"`
}

func (r ticketRow) toDomain() (*domain.Ticket, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, err
	}

	requesterID, err := uuid.Parse(r.RequesterID)
	if err != nil {
		return nil, err
	}

	var assigneeID *uuid.UUID
	if r.AssigneeID != nil {
		parsed, err := uuid.Parse(*r.AssigneeID)
		if err != nil {
			return nil, err
		}
		assigneeID = &parsed
	}

	createdAt, err := time.Parse(time.RFC3339Nano, r.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Ticket{
		ID:          id,
		Title:       r.Title,
		Description: r.Description,
		Status:      domain.Status(r.Status),
		Priority:    domain.Priority(r.Priority),
		RequesterID: requesterID,
		AssigneeID:  assigneeID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func fromDomain(t *domain.Ticket) ticketRow {
	var assigneeID *string
	if t.AssigneeID != nil {
		id := t.AssigneeID.String()
		assigneeID = &id
	}

	return ticketRow{
		ID:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		RequesterID: t.RequesterID.String(),
		AssigneeID:  assigneeID,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339Nano),
	}
}
