package handler_test

import (
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

// TestUpdateTicket_DoSPrevention verifies that authorization checks happen BEFORE body parsing.
// If body parsing happens first, a large body triggers 413 Payload Too Large.
// If auth happens first, an unauthorized user gets 404 (Ticket Not Found for them) or 403 Forbidden without reading the body.
func TestUpdateTicket_DoSPrevention(t *testing.T) {
	mockService := new(MockTicketService)
	mockOrgRepo := new(MockOrgRepo)
	// We don't need UserRepo for this test as we inject user in context
	h := handler.NewTicketHandler(mockService, mockOrgRepo, nil, nil)
	r := chi.NewRouter()
	r.Patch("/tickets/{ticketID}", h.UpdateTicket)

	// User belongs to Org A
	user := &domain.User{ID: uuid.New(), Role: domain.RoleStaff}
	// Ticket belongs to Org B
	orgB := uuid.New()
	ticketID := uuid.New()
	ticket := &domain.Ticket{ID: ticketID, OrganizationID: orgB}

	// Mock Expectation: GetTicket IS called
	mockService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil)
	// Mock Expectation: GetMemberRole IS called and returns "" (unauthorized)
	mockOrgRepo.On("GetMemberRole", mock.Anything, orgB, user.ID).Return("", nil)

	// Create a large body reader (larger than 32MB)
	// MaxRequestSize is 32MB. We send 32MB + 1KB.
	largeBody := &LargeReader{Size: (32 << 20) + 1024}

	req := httptest.NewRequest("PATCH", "/tickets/"+ticketID.String(), largeBody)
	req.Header.Set("Content-Type", "application/json")

	// Inject User
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// CURRENT BEHAVIOR: 413 Payload Too Large (because body is read before auth)
	// DESIRED BEHAVIOR: 404 Not Found (because auth check fails before reading body)

	// DESIRED BEHAVIOR: 404 Not Found (because auth check fails before reading body)

	// We expect 404 because the handler returns 404 if ticket is not found for the user (auth check).
	// If it returns 413, it means it read the body first.
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected 404 Not Found, confirming auth check happened before body read")
}
