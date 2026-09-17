package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/jhonsferg/corelog/internal/platform/apperror"
	"github.com/jhonsferg/corelog/internal/ticket/application/port"
	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

const pgForeignKeyViolation = "23503"

type TicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

var _ port.Repository = (*TicketRepository)(nil)

func (r *TicketRepository) Save(ctx context.Context, t *domain.Ticket) error {
	row := fromDomain(t)

	const query = `
		INSERT INTO tickets (id, title, description, status, priority, requester_id, assignee_id, created_at, updated_at)
		VALUES (:id, :title, :description, :status, :priority, :requester_id, :assignee_id, :created_at, :updated_at)`

	if _, err := r.db.NamedExecContext(ctx, query, row); err != nil {
		if isForeignKeyViolation(err) {
			return domain.ErrRequesterNotFound
		}
		return apperror.Wrap(apperror.KindInternal, "failed to save ticket", err)
	}
	return nil
}

func (r *TicketRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	var row ticketRow

	const query = `
		SELECT id, title, description, status, priority, requester_id, assignee_id, created_at, updated_at
		FROM tickets WHERE id = $1`

	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTicketNotFound
		}
		return nil, apperror.Wrap(apperror.KindInternal, "failed to find ticket", err)
	}

	return row.toDomain(), nil
}

func (r *TicketRepository) FindAll(ctx context.Context, filter port.ListFilter) ([]*domain.Ticket, error) {
	query := `
		SELECT id, title, description, status, priority, requester_id, assignee_id, created_at, updated_at
		FROM tickets WHERE 1 = 1`
	args := []any{}

	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf(" AND status = $%d", len(args))
	}
	if filter.Priority != "" {
		args = append(args, filter.Priority)
		query += fmt.Sprintf(" AND priority = $%d", len(args))
	}

	query += " ORDER BY created_at DESC"

	args = append(args, filter.PageSize)
	query += fmt.Sprintf(" LIMIT $%d", len(args))

	args = append(args, (filter.Page-1)*filter.PageSize)
	query += fmt.Sprintf(" OFFSET $%d", len(args))

	var rows []ticketRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, apperror.Wrap(apperror.KindInternal, "failed to list tickets", err)
	}

	tickets := make([]*domain.Ticket, 0, len(rows))
	for _, row := range rows {
		tickets = append(tickets, row.toDomain())
	}
	return tickets, nil
}

func (r *TicketRepository) Update(ctx context.Context, t *domain.Ticket) error {
	row := fromDomain(t)

	const query = `
		UPDATE tickets
		SET title = :title, description = :description, status = :status, priority = :priority,
		    assignee_id = :assignee_id, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, row)
	if err != nil {
		if isForeignKeyViolation(err) {
			return domain.ErrAssigneeNotFound
		}
		return apperror.Wrap(apperror.KindInternal, "failed to update ticket", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return apperror.Wrap(apperror.KindInternal, "failed to update ticket", err)
	}
	if affected == 0 {
		return domain.ErrTicketNotFound
	}

	return nil
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation
}
