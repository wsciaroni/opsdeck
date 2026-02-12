package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
	"github.com/wsciaroni/opsdeck/internal/core/port"
)

// Helper to generate large strings
func RepeatString(s string, count int) string {
	return strings.Repeat(s, count)
}

// MockScheduledTaskService for testing
type MockScheduledTaskService struct {
	mock.Mock
}

func (m *MockScheduledTaskService) CreateTask(ctx context.Context, cmd port.CreateScheduledTaskCmd) (*domain.ScheduledTask, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduledTask), args.Error(1)
}

func (m *MockScheduledTaskService) UpdateTask(ctx context.Context, id uuid.UUID, cmd port.UpdateScheduledTaskCmd) (*domain.ScheduledTask, error) {
	args := m.Called(ctx, id, cmd)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduledTask), args.Error(1)
}

func (m *MockScheduledTaskService) GetTask(ctx context.Context, id uuid.UUID) (*domain.ScheduledTask, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScheduledTask), args.Error(1)
}

func (m *MockScheduledTaskService) ListTasks(ctx context.Context, orgID uuid.UUID) ([]domain.ScheduledTask, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduledTask), args.Error(1)
}

func (m *MockScheduledTaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func TestOrgHandler_Security(t *testing.T) {
	mockOrgRepo := new(MockOrgRepo)
	mockUserRepo := new(MockUserRepo)
	mockNotificationService := new(MockNotificationService)
	h := handler.NewOrgHandler(mockOrgRepo, mockUserRepo, mockNotificationService, nil)

	r := chi.NewRouter()
	r.Post("/organizations", h.CreateOrganization)

	t.Run("DoS Prevention - CreateOrganization Body Too Large", func(t *testing.T) {
		// Create a body larger than 1MB with valid JSON structure start
		// {"name": "AAAA..."}
		padding := RepeatString("A", 1048576)
		reqBody := map[string]string{"name": padding}
		bodyBytes, _ := json.Marshal(reqBody)
		// This might be slightly larger than 1MB + overhead.

		req := httptest.NewRequest("POST", "/organizations", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Setup user context
		user := &domain.User{ID: uuid.New()}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	})

	t.Run("Input Validation - Name Too Long", func(t *testing.T) {
		longName := RepeatString("A", 101)
		reqBody := map[string]string{"name": longName}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/organizations", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Setup user context
		user := &domain.User{ID: uuid.New()}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Name too long")
	})
}

func TestScheduledTaskHandler_Security(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Post("/scheduled-tasks", h.Create)
	r.Patch("/scheduled-tasks/{id}", h.Update)

	user := &domain.User{ID: uuid.New()}
	orgID := uuid.New()

	// Pre-setup for membership check which happens AFTER body read/validation in some handlers,
	// but BEFORE in others?
	// In Create: Decode -> User Check -> Membership Check.
	// So validation happens AFTER Decode.
	// If Decode fails due to size, it returns early.
	// If Decode succeeds, it validates fields.
	// So we don't need mocks for Body Too Large or Invalid Input.

	t.Run("DoS Prevention - Create Task Body Too Large", func(t *testing.T) {
		padding := RepeatString("A", 1048576)
		reqBody := map[string]string{"title": padding}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/scheduled-tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	})

	t.Run("Input Validation - Title Too Long", func(t *testing.T) {
		longTitle := RepeatString("A", 201)
		reqBody := map[string]interface{}{
			"organization_id": orgID,
			"title": longTitle,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/scheduled-tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Title must be between")
	})

	t.Run("Input Validation - Description Too Long", func(t *testing.T) {
		longDesc := RepeatString("A", 5001)
		reqBody := map[string]interface{}{
			"organization_id": orgID,
			"title": "Valid Title",
			"description": longDesc,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/scheduled-tasks", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Description too long")
	})

	t.Run("DoS Prevention - Update Task Body Too Large", func(t *testing.T) {
		// Expect GetTask and Auth check because we moved them before body read
		taskID := uuid.New()
		task := &domain.ScheduledTask{ID: taskID, OrganizationID: orgID}
		mockService.On("GetTask", mock.Anything, taskID).Return(task, nil)
		mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("owner", nil)

		padding := RepeatString("A", 1048576)
		reqBody := map[string]string{"title": padding}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PATCH", "/scheduled-tasks/"+taskID.String(), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	})
}
