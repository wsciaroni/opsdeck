package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

type ScheduledTaskHandler struct {
	service port.ScheduledTaskService
	orgRepo port.OrganizationRepository
	logger  *slog.Logger
}

func NewScheduledTaskHandler(service port.ScheduledTaskService, orgRepo port.OrganizationRepository, logger *slog.Logger) *ScheduledTaskHandler {
	return &ScheduledTaskHandler{
		service: service,
		orgRepo: orgRepo,
		logger:  logger,
	}
}

type CreateScheduledTaskRequest struct {
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Frequency      domain.Frequency `json:"frequency"`
	StartDate      time.Time        `json:"start_date"`
	PriorityID     string           `json:"priority_id"`
	OrganizationID uuid.UUID        `json:"organization_id"`
	AssigneeUserID *uuid.UUID       `json:"assignee_user_id"`
	Location       string           `json:"location"`
	Enabled        bool             `json:"enabled"`
}

type UpdateScheduledTaskRequest struct {
	Title          *string           `json:"title"`
	Description    *string           `json:"description"`
	Frequency      *domain.Frequency `json:"frequency"`
	StartDate      *time.Time        `json:"start_date"`
	PriorityID     *string           `json:"priority_id"`
	AssigneeUserID *uuid.UUID        `json:"assignee_user_id"`
	Location       *string           `json:"location"`
	Enabled        *bool             `json:"enabled"`
}

func (h *ScheduledTaskHandler) List(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("organization_id")
	if orgIDStr == "" {
		http.Error(w, "organization_id is required", http.StatusBadRequest)
		return
	}
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.verifyMembership(r.Context(), user.ID, orgID); err != nil {
		h.logger.Error("failed to verify membership", "error", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	tasks, err := h.service.ListTasks(r.Context(), orgID)
	if err != nil {
		h.logger.Error("failed to list tasks", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *ScheduledTaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateScheduledTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.verifyMembership(r.Context(), user.ID, req.OrganizationID); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	cmd := port.CreateScheduledTaskCmd{
		OrganizationID: req.OrganizationID,
		CreatedBy:      user.ID,
		Title:          req.Title,
		Description:    req.Description,
		Frequency:      req.Frequency,
		StartDate:      req.StartDate,
		AssigneeUserID: req.AssigneeUserID,
		PriorityID:     req.PriorityID,
		Location:       req.Location,
		Enabled:        req.Enabled,
	}

	task, err := h.service.CreateTask(r.Context(), cmd)
	if err != nil {
		h.logger.Error("failed to create task", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *ScheduledTaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateScheduledTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.verifyMembership(r.Context(), user.ID, task.OrganizationID); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	cmd := port.UpdateScheduledTaskCmd{
		Title:          req.Title,
		Description:    req.Description,
		Frequency:      req.Frequency,
		StartDate:      req.StartDate,
		AssigneeUserID: req.AssigneeUserID,
		PriorityID:     req.PriorityID,
		Location:       req.Location,
		Enabled:        req.Enabled,
	}

	updated, err := h.service.UpdateTask(r.Context(), id, cmd)
	if err != nil {
		h.logger.Error("failed to update task", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *ScheduledTaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.verifyMembership(r.Context(), user.ID, task.OrganizationID); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		h.logger.Error("failed to delete task", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ScheduledTaskHandler) verifyMembership(ctx context.Context, userID, orgID uuid.UUID) error {
	memberships, err := h.orgRepo.ListByUser(ctx, userID)
	if err != nil {
		return err
	}
	for _, m := range memberships {
		if m.ID == orgID {
			return nil
		}
	}
	return fmt.Errorf("user not member of organization")
}
