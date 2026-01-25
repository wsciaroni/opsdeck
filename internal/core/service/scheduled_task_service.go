package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type ScheduledTaskService struct {
	repo port.ScheduledTaskRepository
}

func NewScheduledTaskService(repo port.ScheduledTaskRepository) *ScheduledTaskService {
	return &ScheduledTaskService{repo: repo}
}

func (s *ScheduledTaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.ScheduledTask, error) {
	return s.repo.Get(ctx, id)
}

func (s *ScheduledTaskService) ListTasks(ctx context.Context, organizationID uuid.UUID) ([]domain.ScheduledTask, error) {
	return s.repo.List(ctx, organizationID)
}

func (s *ScheduledTaskService) CreateTask(ctx context.Context, cmd port.CreateScheduledTaskCmd) (*domain.ScheduledTask, error) {
	nextRun := calculateNextRun(cmd.StartDate, cmd.Frequency)

	task := &domain.ScheduledTask{
		OrganizationID: cmd.OrganizationID,
		Title:          cmd.Title,
		Description:    cmd.Description,
		Frequency:      cmd.Frequency,
		StartDate:      cmd.StartDate,
		NextRunAt:      nextRun,
		CreatedBy:      cmd.CreatedBy,
		AssigneeUserID: cmd.AssigneeUserID,
		PriorityID:     cmd.PriorityID,
		Location:       cmd.Location,
		Enabled:        cmd.Enabled,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *ScheduledTaskService) UpdateTask(ctx context.Context, id uuid.UUID, cmd port.UpdateScheduledTaskCmd) (*domain.ScheduledTask, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("scheduled task not found")
	}

	if cmd.Title != nil {
		task.Title = *cmd.Title
	}
	if cmd.Description != nil {
		task.Description = *cmd.Description
	}
	if cmd.Frequency != nil {
		task.Frequency = *cmd.Frequency
	}
	if cmd.StartDate != nil {
		task.StartDate = *cmd.StartDate
	}
	if cmd.PriorityID != nil {
		task.PriorityID = *cmd.PriorityID
	}
	if cmd.Location != nil {
		task.Location = *cmd.Location
	}
	if cmd.Enabled != nil {
		task.Enabled = *cmd.Enabled
	}
	if cmd.AssigneeUserID != nil {
		if *cmd.AssigneeUserID == uuid.Nil {
			task.AssigneeUserID = nil
		} else {
			task.AssigneeUserID = cmd.AssigneeUserID
		}
	}

	if cmd.Frequency != nil || cmd.StartDate != nil {
		task.NextRunAt = calculateNextRun(task.StartDate, task.Frequency)
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *ScheduledTaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func calculateNextRun(start time.Time, freq domain.Frequency) time.Time {
	now := time.Now()
	next := start

	if next.After(now) {
		return next
	}

	for !next.After(now) {
		switch freq {
		case domain.FrequencyDaily:
			next = next.AddDate(0, 0, 1)
		case domain.FrequencyWeekly:
			next = next.AddDate(0, 0, 7)
		case domain.FrequencyMonthly:
			next = next.AddDate(0, 1, 0)
		case domain.FrequencyYearly:
			next = next.AddDate(1, 0, 0)
		}
	}
	return next
}
