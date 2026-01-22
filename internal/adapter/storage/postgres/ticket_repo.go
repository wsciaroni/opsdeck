package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type TicketRepository struct {
	db *pgxpool.Pool
}

func NewTicketRepository(db *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	if ticket.OrganizationID == uuid.Nil {
		return fmt.Errorf("organization_id cannot be nil")
	}

	query := `
		INSERT INTO tickets (
			organization_id, reporter_id, assignee_user_id, status_id, priority_id,
			title, description, location, completed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		ticket.OrganizationID,
		ticket.ReporterID,
		ticket.AssigneeUserID,
		ticket.StatusID,
		ticket.PriorityID,
		ticket.Title,
		ticket.Description,
		ticket.Location,
		ticket.CompletedAt,
	).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	return nil
}

func (r *TicketRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	query := `
		SELECT id, organization_id, reporter_id, assignee_user_id, status_id, priority_id,
		       title, description, location, created_at, updated_at, completed_at
		FROM tickets
		WHERE id = $1
	`

	var t domain.Ticket
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.OrganizationID,
		&t.ReporterID,
		&t.AssigneeUserID,
		&t.StatusID,
		&t.PriorityID,
		&t.Title,
		&t.Description,
		&t.Location,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.CompletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ticket by id: %w", err)
	}

	return &t, nil
}

func (r *TicketRepository) List(ctx context.Context, filter port.TicketFilter) ([]domain.Ticket, error) {
	descriptionField := "description"
	if filter.ExcludeDescription {
		descriptionField = "'' as description"
	}

	query := fmt.Sprintf(`
		SELECT id, organization_id, reporter_id, assignee_user_id, status_id, priority_id,
		       title, %s, location, created_at, updated_at, completed_at
		FROM tickets
		WHERE 1=1
	`, descriptionField)
	args := []interface{}{}
	argIdx := 1

	if filter.OrganizationID != nil {
		query += fmt.Sprintf(" AND organization_id = $%d", argIdx)
		args = append(args, *filter.OrganizationID)
		argIdx++
	} else if len(filter.OrganizationIDs) > 0 {
		query += fmt.Sprintf(" AND organization_id = ANY($%d)", argIdx)
		args = append(args, filter.OrganizationIDs)
		argIdx++
	}

	if filter.StatusID != nil {
		query += fmt.Sprintf(" AND status_id = $%d", argIdx)
		args = append(args, *filter.StatusID)
		argIdx++
	}

	if filter.AssigneeID != nil {
		query += fmt.Sprintf(" AND assignee_user_id = $%d", argIdx)
		args = append(args, *filter.AssigneeID)
		argIdx++
	}

	if filter.ReporterID != nil {
		query += fmt.Sprintf(" AND reporter_id = $%d", argIdx)
		args = append(args, *filter.ReporterID)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		var t domain.Ticket
		err := rows.Scan(
			&t.ID,
			&t.OrganizationID,
			&t.ReporterID,
			&t.AssigneeUserID,
			&t.StatusID,
			&t.PriorityID,
			&t.Title,
			&t.Description,
			&t.Location,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		tickets = append(tickets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tickets, nil
}

func (r *TicketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	query := `
		UPDATE tickets
		SET status_id = $1, priority_id = $2, assignee_user_id = $3,
		    title = $4, description = $5, location = $6,
		    updated_at = $7, completed_at = $8
		WHERE id = $9
	`

	tag, err := r.db.Exec(ctx, query,
		ticket.StatusID,
		ticket.PriorityID,
		ticket.AssigneeUserID,
		ticket.Title,
		ticket.Description,
		ticket.Location,
		ticket.UpdatedAt,
		ticket.CompletedAt,
		ticket.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ticket not found")
	}

	return nil
}
