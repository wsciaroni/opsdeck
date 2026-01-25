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
			title, description, location, completed_at, sensitive
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		ticket.Sensitive,
	).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	return nil
}

func (r *TicketRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	query := `
		SELECT id, organization_id, reporter_id, assignee_user_id, status_id, priority_id,
		       title, description, location, created_at, updated_at, completed_at, sensitive
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
		&t.Sensitive,
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
		       title, %s, location, created_at, updated_at, completed_at, sensitive
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

	if len(filter.StatusIDs) > 0 {
		query += fmt.Sprintf(" AND status_id = ANY($%d)", argIdx)
		args = append(args, filter.StatusIDs)
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
		argIdx++
	}

	if filter.Sensitive != nil {
		query += fmt.Sprintf(" AND sensitive = $%d", argIdx)
		args = append(args, *filter.Sensitive)
		argIdx++
	}

	if len(filter.PriorityIDs) > 0 {
		query += fmt.Sprintf(" AND priority_id = ANY($%d)", argIdx)
		args = append(args, filter.PriorityIDs)
		argIdx++
	}

	if filter.Keyword != nil && *filter.Keyword != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx)
		keyword := fmt.Sprintf("%%%s%%", *filter.Keyword)
		args = append(args, keyword)
	}

	orderBy := "created_at DESC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}

		switch filter.SortBy {
		case "created_at":
			orderBy = fmt.Sprintf("created_at %s", direction)
		case "updated_at":
			orderBy = fmt.Sprintf("updated_at %s", direction)
		case "title":
			orderBy = fmt.Sprintf("title %s", direction)
		case "priority":
			orderBy = fmt.Sprintf(`CASE priority_id
				WHEN 'critical' THEN 4
				WHEN 'high' THEN 3
				WHEN 'medium' THEN 2
				WHEN 'low' THEN 1
				ELSE 0
			END %s`, direction)
		case "status":
			orderBy = fmt.Sprintf(`CASE status_id
				WHEN 'new' THEN 1
				WHEN 'in_progress' THEN 2
				WHEN 'on_hold' THEN 3
				WHEN 'done' THEN 4
				WHEN 'canceled' THEN 5
				ELSE 6
			END %s`, direction)
		}
	}

	query += " ORDER BY " + orderBy

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
			&t.Sensitive,
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
		    updated_at = $7, completed_at = $8, sensitive = $9
		WHERE id = $10
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
		ticket.Sensitive,
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

func (r *TicketRepository) AddFile(ctx context.Context, file *domain.File) error {
	query := `
		INSERT INTO ticket_files (ticket_id, filename, content_type, size, data)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query,
		file.TicketID,
		file.Filename,
		file.ContentType,
		file.Size,
		file.Data,
	).Scan(&file.ID, &file.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to add file: %w", err)
	}
	return nil
}

func (r *TicketRepository) GetFile(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	query := `
		SELECT id, ticket_id, filename, content_type, size, data, created_at
		FROM ticket_files
		WHERE id = $1
	`
	var f domain.File
	err := r.db.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.TicketID,
		&f.Filename,
		&f.ContentType,
		&f.Size,
		&f.Data,
		&f.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return &f, nil
}

func (r *TicketRepository) ListFiles(ctx context.Context, ticketID uuid.UUID) ([]domain.File, error) {
	query := `
		SELECT id, ticket_id, filename, content_type, size, created_at
		FROM ticket_files
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`
	// Note: NOT selecting data.

	rows, err := r.db.Query(ctx, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer rows.Close()

	var files []domain.File
	for rows.Next() {
		var f domain.File
		err := rows.Scan(
			&f.ID,
			&f.TicketID,
			&f.Filename,
			&f.ContentType,
			&f.Size,
			&f.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, f)
	}
	return files, nil
}
