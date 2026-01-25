package port

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

type CreateScheduledTaskCmd struct {
	OrganizationID uuid.UUID
	CreatedBy      uuid.UUID
	Title          string
	Description    string
	Frequency      domain.Frequency
	StartDate      time.Time
	AssigneeUserID *uuid.UUID
	PriorityID     string
	Location       string
	Enabled        bool
}

type UpdateScheduledTaskCmd struct {
	Title          *string
	Description    *string
	Frequency      *domain.Frequency
	StartDate      *time.Time
	AssigneeUserID *uuid.UUID
	PriorityID     *string
	Location       *string
	Enabled        *bool
}

type ScheduledTaskService interface {
	GetTask(ctx context.Context, id uuid.UUID) (*domain.ScheduledTask, error)
	ListTasks(ctx context.Context, organizationID uuid.UUID) ([]domain.ScheduledTask, error)
	CreateTask(ctx context.Context, cmd CreateScheduledTaskCmd) (*domain.ScheduledTask, error)
	UpdateTask(ctx context.Context, id uuid.UUID, cmd UpdateScheduledTaskCmd) (*domain.ScheduledTask, error)
	DeleteTask(ctx context.Context, id uuid.UUID) error
}
