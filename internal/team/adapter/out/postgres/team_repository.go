package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/jhonsferg/corelog/internal/platform/apperror"
	"github.com/jhonsferg/corelog/internal/team/application/port"
	"github.com/jhonsferg/corelog/internal/team/domain"
)

type TeamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

var _ port.Repository = (*TeamRepository)(nil)

func (r *TeamRepository) Save(ctx context.Context, t *domain.Team) error {
	row := fromDomain(t)

	const query = `
		INSERT INTO teams (id, name, created_at, updated_at)
		VALUES (:id, :name, :created_at, :updated_at)`

	if _, err := r.db.NamedExecContext(ctx, query, row); err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to save team", err)
	}
	return nil
}

func (r *TeamRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	var row teamRow

	const query = `SELECT id, name, created_at, updated_at FROM teams WHERE id = $1`

	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, apperror.Wrap(apperror.KindInternal, "failed to find team", err)
	}

	return row.toDomain(), nil
}

func (r *TeamRepository) FindByName(ctx context.Context, name string) (*domain.Team, error) {
	var row teamRow

	const query = `SELECT id, name, created_at, updated_at FROM teams WHERE name = $1`

	if err := r.db.GetContext(ctx, &row, query, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, apperror.Wrap(apperror.KindInternal, "failed to find team", err)
	}

	return row.toDomain(), nil
}

func (r *TeamRepository) FindAll(ctx context.Context) ([]*domain.Team, error) {
	var rows []teamRow

	const query = `SELECT id, name, created_at, updated_at FROM teams ORDER BY name ASC`

	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, apperror.Wrap(apperror.KindInternal, "failed to list teams", err)
	}

	teams := make([]*domain.Team, 0, len(rows))
	for _, row := range rows {
		teams = append(teams, row.toDomain())
	}
	return teams, nil
}

func (r *TeamRepository) Update(ctx context.Context, t *domain.Team) error {
	row := fromDomain(t)

	const query = `UPDATE teams SET name = :name, updated_at = :updated_at WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to update team", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to update team", err)
	}
	if affected == 0 {
		return domain.ErrTeamNotFound
	}

	return nil
}

func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM teams WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to delete team", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to delete team", err)
	}
	if affected == 0 {
		return domain.ErrTeamNotFound
	}

	return nil
}
