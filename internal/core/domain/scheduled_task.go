package domain

import (
	"time"

	"github.com/google/uuid"
)

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
	FrequencyYearly  Frequency = "yearly"
)

type ScheduledTask struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Frequency      Frequency  `json:"frequency"`
	StartDate      time.Time  `json:"start_date"`
	NextRunAt      time.Time  `json:"next_run_at"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	AssigneeUserID *uuid.UUID `json:"assignee_user_id"`
	PriorityID     string     `json:"priority_id"`
	Location       string     `json:"location"`
	Enabled        bool       `json:"enabled"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
