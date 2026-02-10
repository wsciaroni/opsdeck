package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

// MockScheduledTaskService is defined in security_test.go (same package)

func TestCreateTask_Forbidden_Member(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Post("/scheduled-tasks", h.Create)

	user := &domain.User{ID: uuid.New(), Role: domain.RoleStaff}
	orgID := uuid.New()

	// Mock GetMemberRole returning "member" (not owner/admin)
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil)

	reqBody := map[string]interface{}{
		"title":           "New Task",
		"description":     "Do something",
		"organization_id": orgID,
		"frequency":       "daily",
		"start_date":      time.Now(),
		"priority_id":     "medium",
		"enabled":         true,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/scheduled-tasks", bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Expect Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateTask_Forbidden_Member(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Patch("/scheduled-tasks/{id}", h.Update)

	user := &domain.User{ID: uuid.New(), Role: domain.RoleStaff}
	orgID := uuid.New()
	taskID := uuid.New()
	task := &domain.ScheduledTask{ID: taskID, OrganizationID: orgID}

	mockService.On("GetTask", mock.Anything, taskID).Return(task, nil)
	// Mock GetMemberRole returning "member"
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil)

	reqBody := map[string]interface{}{
		"title": "Updated Title",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PATCH", "/scheduled-tasks/"+taskID.String(), bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Expect Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteTask_Forbidden_Member(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Delete("/scheduled-tasks/{id}", h.Delete)

	user := &domain.User{ID: uuid.New(), Role: domain.RoleStaff}
	orgID := uuid.New()
	taskID := uuid.New()
	task := &domain.ScheduledTask{ID: taskID, OrganizationID: orgID}

	mockService.On("GetTask", mock.Anything, taskID).Return(task, nil)
	// Mock GetMemberRole returning "member"
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil)

	req := httptest.NewRequest("DELETE", "/scheduled-tasks/"+taskID.String(), nil)
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Expect Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteTask_Success_Admin(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Delete("/scheduled-tasks/{id}", h.Delete)

	user := &domain.User{ID: uuid.New(), Role: domain.RoleAdmin}
	orgID := uuid.New()
	taskID := uuid.New()
	task := &domain.ScheduledTask{ID: taskID, OrganizationID: orgID}

	mockService.On("GetTask", mock.Anything, taskID).Return(task, nil)
	// Mock GetMemberRole returning "admin"
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("admin", nil)
	mockService.On("DeleteTask", mock.Anything, taskID).Return(nil)

	req := httptest.NewRequest("DELETE", "/scheduled-tasks/"+taskID.String(), nil)
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Expect No Content
	assert.Equal(t, http.StatusNoContent, w.Code)
}
