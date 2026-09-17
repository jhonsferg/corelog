package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/jhonsferg/corelog/internal/user/domain"
)

type Repository interface {
	Save(ctx context.Context, u *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Search(ctx context.Context, filter SearchFilter) ([]*domain.User, error)
	Count(ctx context.Context) (int, error)
}
