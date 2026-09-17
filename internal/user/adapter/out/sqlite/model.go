package sqlite

import (
	"time"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/user/domain"
)

type userRow struct {
	ID           string  `db:"id"`
	Name         string  `db:"name"`
	Email        string  `db:"email"`
	PasswordHash string  `db:"password_hash"`
	Role         string  `db:"role"`
	TeamID       *string `db:"team_id"`
	CreatedAt    string  `db:"created_at"`
	UpdatedAt    string  `db:"updated_at"`
}

func (r userRow) toDomain() (*domain.User, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, err
	}

	var teamID *uuid.UUID
	if r.TeamID != nil {
		parsed, err := uuid.Parse(*r.TeamID)
		if err != nil {
			return nil, err
		}
		teamID = &parsed
	}

	createdAt, err := time.Parse(time.RFC3339Nano, r.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           id,
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		Role:         domain.Role(r.Role),
		TeamID:       teamID,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func fromDomain(u *domain.User) userRow {
	var teamID *string
	if u.TeamID != nil {
		id := u.TeamID.String()
		teamID = &id
	}

	return userRow{
		ID:           u.ID.String(),
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		TeamID:       teamID,
		CreatedAt:    u.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339Nano),
	}
}
