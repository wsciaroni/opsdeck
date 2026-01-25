package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

type ScheduledTaskRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.ScheduledTask, error)
	List(ctx context.Context, organizationID uuid.UUID) ([]domain.ScheduledTask, error)
	Create(ctx context.Context, task *domain.ScheduledTask) error
	Update(ctx context.Context, task *domain.ScheduledTask) error
	Delete(ctx context.Context, id uuid.UUID) error
}
