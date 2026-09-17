package http

import (
	"time"

	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

type createTicketRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"required,min=3"`
	Priority    string `json:"priority" validate:"required,oneof=low medium high critical"`
}

type updateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=open in_progress resolved closed"`
}

type assignRequest struct {
	AssigneeID string `json:"assignee_id" validate:"required,uuid4"`
}

type ticketResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	RequesterID string    `json:"requester_id"`
	AssigneeID  *string   `json:"assignee_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toTicketResponse(t *domain.Ticket) ticketResponse {
	resp := ticketResponse{
		ID:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		RequesterID: t.RequesterID.String(),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
	if t.AssigneeID != nil {
		id := t.AssigneeID.String()
		resp.AssigneeID = &id
	}
	return resp
}

func toTicketResponses(tickets []*domain.Ticket) []ticketResponse {
	out := make([]ticketResponse, 0, len(tickets))
	for _, t := range tickets {
		out = append(out, toTicketResponse(t))
	}
	return out
}
