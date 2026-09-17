package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/jhonsferg/corelog/internal/platform/apperror"
	"github.com/jhonsferg/corelog/internal/ticket/application/port"
	"github.com/jhonsferg/corelog/internal/ticket/domain"
)

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
		FROM tickets WHERE id = ?`

	if err := r.db.GetContext(ctx, &row, query, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTicketNotFound
		}
		return nil, apperror.Wrap(apperror.KindInternal, "failed to find ticket", err)
	}

	ticket, err := row.toDomain()
	if err != nil {
		return nil, apperror.Wrap(apperror.KindInternal, "failed to decode ticket", err)
	}
	return ticket, nil
}

func (r *TicketRepository) FindAll(ctx context.Context, filter port.ListFilter) ([]*domain.Ticket, error) {
	query := `
		SELECT id, title, description, status, priority, requester_id, assignee_id, created_at, updated_at
		FROM tickets WHERE 1 = 1`
	args := []any{}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, string(filter.Status))
	}
	if filter.Priority != "" {
		query += " AND priority = ?"
		args = append(args, string(filter.Priority))
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	var rows []ticketRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, apperror.Wrap(apperror.KindInternal, "failed to list tickets", err)
	}

	tickets := make([]*domain.Ticket, 0, len(rows))
	for _, row := range rows {
		ticket, err := row.toDomain()
		if err != nil {
			return nil, apperror.Wrap(apperror.KindInternal, "failed to decode ticket", err)
		}
		tickets = append(tickets, ticket)
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
	return strings.Contains(err.Error(), "FOREIGN KEY constraint failed")
}
