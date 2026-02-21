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

func TestCreateTask_Optimization_SingleDBCheck(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Post("/scheduled-tasks", h.Create)

	user := &domain.User{ID: uuid.New(), Role: domain.RoleAdmin}
	orgID := uuid.New()

	// Mock Expectation: GetMemberRole should be called ONCE
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("admin", nil).Once()

	mockService.On("CreateTask", mock.Anything, mock.AnythingOfType("port.CreateScheduledTaskCmd")).Return(&domain.ScheduledTask{ID: uuid.New()}, nil)

	reqBody := map[string]interface{}{
		"title":           "Optimized Task",
		"description":     "Should only check auth once",
		"organization_id": orgID,
		"frequency":       "daily",
		"start_date":      time.Now(),
		"priority_id":     "medium",
		"enabled":         true,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/scheduled-tasks?organization_id="+orgID.String(), bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockOrgRepo.AssertExpectations(t)
}

func TestCreateTask_Optimization_ExploitAttempt(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Post("/scheduled-tasks", h.Create)

	user := &domain.User{ID: uuid.New(), Role: domain.RoleAdmin}
	orgID_A := uuid.New() // Authorized
	orgID_B := uuid.New() // Unauthorized

	// User is admin of A
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID_A, user.ID).Return("admin", nil).Once()

	// User is NOT checked for B (because mismatch check happens first)
	// But if code is buggy and proceeds to use B without check, it would call GetMemberRole for B (and fail or succeed if mocked)
	// We assert that GetMemberRole is NOT called for B.
	// But since we use strict mocks, any unexpected call fails.

	reqBody := map[string]interface{}{
		"title":           "Exploit Task",
		"description":     "Trying to create task in Org B using auth for Org A",
		"organization_id": orgID_B,
		"frequency":       "daily",
		"start_date":      time.Now(),
		"priority_id":     "medium",
		"enabled":         true,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Send Query Param for A (Authorized), Body for B (Unauthorized)
	req := httptest.NewRequest("POST", "/scheduled-tasks?organization_id="+orgID_A.String(), bytes.NewReader(bodyBytes))
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Expect 400 Bad Request (Mismatch)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Ensure GetMemberRole was called for A (once) and NOT for B
	mockOrgRepo.AssertExpectations(t)
}
