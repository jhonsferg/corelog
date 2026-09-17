package http

import (
	"time"

	"github.com/jhonsferg/corelog/internal/team/domain"
)

type createTeamRequest struct {
	Name string `json:"name" validate:"required,min=2,max=120"`
}

type renameTeamRequest struct {
	Name string `json:"name" validate:"required,min=2,max=120"`
}

type teamResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toTeamResponse(t *domain.Team) teamResponse {
	return teamResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func toTeamResponses(teams []*domain.Team) []teamResponse {
	out := make([]teamResponse, 0, len(teams))
	for _, t := range teams {
		out = append(out, toTeamResponse(t))
	}
	return out
}
