package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/handler"
	"github.com/wsciaroni/opsdeck/internal/adapter/web/middleware"
	"github.com/wsciaroni/opsdeck/internal/core/domain"
)

func TestCreateTask_DoS_Prevention(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Post("/scheduled-tasks", h.Create)

	user := &domain.User{ID: uuid.New()}
	orgID := uuid.New()

	// Mock GetMemberRole returning "" (not a member)
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("", nil)

	// Send request with organization_id query param
	// The body should NOT be read/parsed if auth fails early.
	// We send invalid JSON to prove body is not parsed (if it returns 403 instead of 400).
	reqBody := []byte(`invalid json`)

	req := httptest.NewRequest("POST", "/scheduled-tasks?organization_id="+orgID.String(), bytes.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// If early auth works, we expect 403 Forbidden.
	// If it falls through to body decode, it would be 400 Bad Request (invalid JSON).
	// NOTE: This test will fail until the implementation is updated.
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreateTask_Mismatch_Prevention(t *testing.T) {
	mockService := new(MockScheduledTaskService)
	mockOrgRepo := new(MockOrgRepo)
	h := handler.NewScheduledTaskHandler(mockService, mockOrgRepo, nil)

	r := chi.NewRouter()
	r.Post("/scheduled-tasks", h.Create)

	user := &domain.User{ID: uuid.New()}
	orgID := uuid.New()
	otherOrgID := uuid.New()

	// User IS admin of query param org
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("admin", nil)
	// mockOrgRepo.On("GetMemberRole", mock.Anything, otherOrgID, user.ID).Return("admin", nil) // Not needed if implementation checks query param first?
	// But later implementation might check body org ID again. It should be skipped if checked early?
	// Wait, if early check passes, we proceed to body decode.
	// Then we parse body. If body orgID != query orgID -> 400.
	// If matching, we proceed. Do we check auth again?
	// Implementation detail: we should probably not check auth TWICE if orgID matches.
	// But for this test, we expect 400 Bad Request.

	// Valid JSON body with DIFFERENT org ID
	reqBody := []byte(`{"organization_id": "` + otherOrgID.String() + `", "title": "Test"}`)

	req := httptest.NewRequest("POST", "/scheduled-tasks?organization_id="+orgID.String(), bytes.NewReader(reqBody))
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
