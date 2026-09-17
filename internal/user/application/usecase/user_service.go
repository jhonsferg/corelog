package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jhonsferg/corelog/internal/platform/apperror"
	"github.com/jhonsferg/corelog/internal/user/application/port"
	"github.com/jhonsferg/corelog/internal/user/domain"
)

type UserService struct {
	repo port.Repository
}

func NewUserService(repo port.Repository) *UserService {
	return &UserService{repo: repo}
}

var _ port.Service = (*UserService)(nil)

func (s *UserService) Register(ctx context.Context, in port.RegisterInput) (*domain.User, error) {
	if _, err := s.repo.FindByEmail(ctx, in.Email); err == nil {
		return nil, domain.ErrEmailAlreadyTaken
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Wrap(apperror.KindInternal, "failed to hash password", err)
	}

	role := in.Role
	if role == "" {
		role = domain.RoleRequester
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		role = domain.RoleAdmin
	}

	user, err := domain.NewUser(in.Name, in.Email, string(hash), role)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, in port.AuthenticateInput) (*domain.User, error) {
	user, err := s.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, in port.UpdateProfileInput) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if in.Email != user.Email {
		if existing, err := s.repo.FindByEmail(ctx, in.Email); err == nil && existing.ID != userID {
			return nil, domain.ErrEmailAlreadyTaken
		} else if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
	}

	if err := user.UpdateProfile(in.Name, in.Email); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, in port.ChangePasswordInput) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.CurrentPassword)); err != nil {
		return domain.ErrIncorrectPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to hash password", err)
	}

	user.ChangePasswordHash(string(hash))

	return s.repo.Update(ctx, user)
}

func (s *UserService) SetTeam(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.AssignTeam(teamID)

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) SearchTeammates(ctx context.Context, actorID uuid.UUID, query string) ([]*domain.User, error) {
	actor, err := s.repo.FindByID(ctx, actorID)
	if err != nil {
		return nil, err
	}

	if actor.TeamID == nil {
		return []*domain.User{}, nil
	}

	return s.repo.Search(ctx, port.SearchFilter{Query: query, TeamID: actor.TeamID})
}

func (s *UserService) Search(ctx context.Context, filter port.SearchFilter) ([]*domain.User, error) {
	return s.repo.Search(ctx, filter)
}
