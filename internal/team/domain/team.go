package domain

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTeam(name string) (*Team, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	now := time.Now().UTC()
	return &Team{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (t *Team) Rename(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	t.Name = name
	t.UpdatedAt = time.Now().UTC()
	return nil
}
