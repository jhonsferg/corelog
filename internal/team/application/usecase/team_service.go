package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/team/application/port"
	"github.com/jhonsferg/corelog/internal/team/domain"
)

type TeamService struct {
	repo port.Repository
}

func NewTeamService(repo port.Repository) *TeamService {
	return &TeamService{repo: repo}
}

var _ port.Service = (*TeamService)(nil)

func (s *TeamService) Create(ctx context.Context, name string) (*domain.Team, error) {
	if _, err := s.repo.FindByName(ctx, name); err == nil {
		return nil, domain.ErrTeamNameTaken
	} else if !errors.Is(err, domain.ErrTeamNotFound) {
		return nil, err
	}

	team, err := domain.NewTeam(name)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) Get(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *TeamService) List(ctx context.Context) ([]*domain.Team, error) {
	return s.repo.FindAll(ctx)
}

func (s *TeamService) Rename(ctx context.Context, id uuid.UUID, name string) (*domain.Team, error) {
	team, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing, err := s.repo.FindByName(ctx, name); err == nil && existing.ID != id {
		return nil, domain.ErrTeamNameTaken
	} else if err != nil && !errors.Is(err, domain.ErrTeamNotFound) {
		return nil, err
	}

	if err := team.Rename(name); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
