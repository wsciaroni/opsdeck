package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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
	// Limit request size to 1MB to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, MaxJSONBodySize)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Title) == 0 || len(req.Title) > 200 {
		http.Error(w, "Title must be between 1 and 200 characters", http.StatusBadRequest)
		return
	}
	if len(req.Description) > 5000 {
		http.Error(w, "Description too long (max 5000 chars)", http.StatusBadRequest)
		return
	}
	if len(req.Location) > 200 {
		http.Error(w, "Location too long", http.StatusBadRequest)
		return
	}

	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Security Check: Only Admin/Owner can create tasks
	if err := h.checkAdminOrOwner(r.Context(), user.ID, req.OrganizationID); err != nil {
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

	// Security Check: Get existing task and verify user permissions BEFORE reading body
	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get task", "error", err)
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

	// Security Check: Only Admin/Owner can update tasks
	if err := h.checkAdminOrOwner(r.Context(), user.ID, task.OrganizationID); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdateScheduledTaskRequest
	// Limit request size to 1MB to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, MaxJSONBodySize)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title != nil && (len(*req.Title) == 0 || len(*req.Title) > 200) {
		http.Error(w, "Title must be between 1 and 200 characters", http.StatusBadRequest)
		return
	}
	if req.Description != nil && len(*req.Description) > 5000 {
		http.Error(w, "Description too long (max 5000 chars)", http.StatusBadRequest)
		return
	}
	if req.Location != nil && len(*req.Location) > 200 {
		http.Error(w, "Location too long", http.StatusBadRequest)
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

	// Security Check: Only Admin/Owner can delete tasks
	if err := h.checkAdminOrOwner(r.Context(), user.ID, task.OrganizationID); err != nil {
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
	// Optimize: Use direct role lookup
	role, err := h.orgRepo.GetMemberRole(ctx, orgID, userID)
	if err != nil {
		return err
	}
	if role == "" {
		return fmt.Errorf("user not member of organization")
	}
	return nil
}

func (h *ScheduledTaskHandler) checkAdminOrOwner(ctx context.Context, userID, orgID uuid.UUID) error {
	role, err := h.orgRepo.GetMemberRole(ctx, orgID, userID)
	if err != nil {
		return err
	}
	if role == "owner" || role == "admin" {
		return nil
	}
	return fmt.Errorf("forbidden: requires admin or owner role")
}
