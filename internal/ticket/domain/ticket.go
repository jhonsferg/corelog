package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID          uuid.UUID
	Title       string
	Description string
	Status      Status
	Priority    Priority
	RequesterID uuid.UUID
	AssigneeID  *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTicket(title, description string, priority Priority, requesterID uuid.UUID) (*Ticket, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if !priority.IsValid() {
		return nil, ErrInvalidPriority
	}
	if requesterID == uuid.Nil {
		return nil, ErrInvalidRequester
	}

	now := time.Now().UTC()
	return &Ticket{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      StatusOpen,
		Priority:    priority,
		RequesterID: requesterID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (t *Ticket) ChangeStatus(newStatus Status) error {
	if !newStatus.IsValid() {
		return ErrInvalidStatus
	}
	if !t.Status.CanTransitionTo(newStatus) {
		return ErrIllegalTransition
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Ticket) AssignTo(assigneeID uuid.UUID) error {
	t.AssigneeID = &assigneeID
	if t.Status == StatusOpen {
		if err := t.ChangeStatus(StatusInProgress); err != nil {
			return err
		}
	}
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Ticket) IsClosed() bool {
	return t.Status == StatusClosed
}
