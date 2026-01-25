package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

type ScheduledTaskRepository struct {
	db *pgxpool.Pool
}

func NewScheduledTaskRepository(db *pgxpool.Pool) *ScheduledTaskRepository {
	return &ScheduledTaskRepository{db: db}
}

func (r *ScheduledTaskRepository) Get(ctx context.Context, id uuid.UUID) (*domain.ScheduledTask, error) {
	query := `
		SELECT id, organization_id, title, description, frequency, start_date, next_run_at,
		       created_by, assignee_user_id, priority_id, location, enabled, created_at, updated_at
		FROM scheduled_tasks
		WHERE id = $1
	`
	var t domain.ScheduledTask
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.OrganizationID,
		&t.Title,
		&t.Description,
		&t.Frequency,
		&t.StartDate,
		&t.NextRunAt,
		&t.CreatedBy,
		&t.AssigneeUserID,
		&t.PriorityID,
		&t.Location,
		&t.Enabled,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get scheduled task: %w", err)
	}
	return &t, nil
}

func (r *ScheduledTaskRepository) List(ctx context.Context, organizationID uuid.UUID) ([]domain.ScheduledTask, error) {
	query := `
		SELECT id, organization_id, title, description, frequency, start_date, next_run_at,
		       created_by, assignee_user_id, priority_id, location, enabled, created_at, updated_at
		FROM scheduled_tasks
		WHERE organization_id = $1
		ORDER BY next_run_at ASC
	`
	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list scheduled tasks: %w", err)
	}
	defer rows.Close()

	var tasks []domain.ScheduledTask
	for rows.Next() {
		var t domain.ScheduledTask
		err := rows.Scan(
			&t.ID,
			&t.OrganizationID,
			&t.Title,
			&t.Description,
			&t.Frequency,
			&t.StartDate,
			&t.NextRunAt,
			&t.CreatedBy,
			&t.AssigneeUserID,
			&t.PriorityID,
			&t.Location,
			&t.Enabled,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scheduled task: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *ScheduledTaskRepository) Create(ctx context.Context, task *domain.ScheduledTask) error {
	query := `
		INSERT INTO scheduled_tasks (
			organization_id, title, description, frequency, start_date, next_run_at,
			created_by, assignee_user_id, priority_id, location, enabled
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		task.OrganizationID,
		task.Title,
		task.Description,
		task.Frequency,
		task.StartDate,
		task.NextRunAt,
		task.CreatedBy,
		task.AssigneeUserID,
		task.PriorityID,
		task.Location,
		task.Enabled,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create scheduled task: %w", err)
	}
	return nil
}

func (r *ScheduledTaskRepository) Update(ctx context.Context, task *domain.ScheduledTask) error {
	query := `
		UPDATE scheduled_tasks
		SET title = $1, description = $2, frequency = $3, start_date = $4, next_run_at = $5,
		    assignee_user_id = $6, priority_id = $7, location = $8, enabled = $9, updated_at = NOW()
		WHERE id = $10
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Frequency,
		task.StartDate,
		task.NextRunAt,
		task.AssigneeUserID,
		task.PriorityID,
		task.Location,
		task.Enabled,
		task.ID,
	).Scan(&task.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("scheduled task not found")
		}
		return fmt.Errorf("failed to update scheduled task: %w", err)
	}
	return nil
}

func (r *ScheduledTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM scheduled_tasks WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete scheduled task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("scheduled task not found")
	}
	return nil
}
