package postgres

import (
	"time"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/user/domain"
)

type userRow struct {
	ID           uuid.UUID  `db:"id"`
	Name         string     `db:"name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Role         string     `db:"role"`
	TeamID       *uuid.UUID `db:"team_id"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

func (r userRow) toDomain() *domain.User {
	return &domain.User{
		ID:           r.ID,
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		Role:         domain.Role(r.Role),
		TeamID:       r.TeamID,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func fromDomain(u *domain.User) userRow {
	return userRow{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		TeamID:       u.TeamID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
