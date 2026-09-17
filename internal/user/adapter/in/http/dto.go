package http

import (
	"time"

	"github.com/jhonsferg/corelog/internal/user/domain"
)

type registerRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=120"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type updateProfileRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=120"`
	Email string `json:"email" validate:"required,email"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type setTeamRequest struct {
	TeamID *string `json:"team_id" validate:"omitempty,uuid4"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TeamID    *string   `json:"team_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func toUserResponse(u *domain.User) userResponse {
	resp := userResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
	}
	if u.TeamID != nil {
		id := u.TeamID.String()
		resp.TeamID = &id
	}
	return resp
}

func toUserResponses(users []*domain.User) []userResponse {
	out := make([]userResponse, 0, len(users))
	for _, u := range users {
		out = append(out, toUserResponse(u))
	}
	return out
}
