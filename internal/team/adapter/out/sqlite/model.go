package sqlite

import (
	"time"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/team/domain"
)

type teamRow struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (r teamRow) toDomain() (*domain.Team, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339Nano, r.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Team{
		ID:        id,
		Name:      r.Name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func fromDomain(t *domain.Team) teamRow {
	return teamRow{
		ID:        t.ID.String(),
		Name:      t.Name,
		CreatedAt: t.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339Nano),
	}
}
