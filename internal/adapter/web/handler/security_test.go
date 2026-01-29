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
)

// LargeReaderString generates a large string reader
func LargeString(size int) string {
	return strings.Repeat("A", size)
}

func TestUpdateTicket_Security(t *testing.T) {
	t.Run("DoS Prevention - Rejects body too large", func(t *testing.T) {
		mockService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		mockUserRepo := new(MockUserRepo)
		h := handler.NewTicketHandler(mockService, mockOrgRepo, mockUserRepo, nil)
		r := chi.NewRouter()
		r.Patch("/tickets/{ticketID}", h.UpdateTicket)

		// Create a reader larger than MaxRequestSize (32MB)
		largeReader := &LargeReader{Size: handler.MaxRequestSize + 1024}

		req := httptest.NewRequest("PATCH", "/tickets/"+uuid.NewString(), largeReader)
		req.Header.Set("Content-Type", "application/json")

		// Authenticated User
		user := &domain.User{ID: uuid.New()}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Expect 413 Payload Too Large
		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	})

	t.Run("DoS Prevention - Rejects title too long", func(t *testing.T) {
		mockService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		h := handler.NewTicketHandler(mockService, mockOrgRepo, nil, nil)
		r := chi.NewRouter()
		r.Patch("/tickets/{ticketID}", h.UpdateTicket)

		ticketID := uuid.New()
		orgID := uuid.New()
		user := &domain.User{ID: uuid.New()}

		// Setup Mocks for successful auth/membership check
		ticket := &domain.Ticket{ID: ticketID, OrganizationID: orgID}
		mockService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil)

		memberships := []domain.UserMembership{{Organization: domain.Organization{ID: orgID}}}
		mockOrgRepo.On("ListByUser", mock.Anything, user.ID).Return(memberships, nil)

		// Payload with massive title
		longTitle := LargeString(250) // Limit is 200
		reqBody := map[string]string{
			"title": longTitle,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PATCH", "/tickets/"+ticketID.String(), bytes.NewReader(bodyBytes))
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Title")
	})

	t.Run("DoS Prevention - Rejects description too long", func(t *testing.T) {
		mockService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		h := handler.NewTicketHandler(mockService, mockOrgRepo, nil, nil)
		r := chi.NewRouter()
		r.Patch("/tickets/{ticketID}", h.UpdateTicket)

		ticketID := uuid.New()
		orgID := uuid.New()
		user := &domain.User{ID: uuid.New()}

		// Setup Mocks
		ticket := &domain.Ticket{ID: ticketID, OrganizationID: orgID}
		mockService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil)

		memberships := []domain.UserMembership{{Organization: domain.Organization{ID: orgID}}}
		mockOrgRepo.On("ListByUser", mock.Anything, user.ID).Return(memberships, nil)

		// Payload with massive description
		longDesc := LargeString(5001) // Limit is 5000
		reqBody := map[string]string{
			"description": longDesc,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PATCH", "/tickets/"+ticketID.String(), bytes.NewReader(bodyBytes))
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Description")
	})

	t.Run("Success - Valid Update", func(t *testing.T) {
		mockService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		h := handler.NewTicketHandler(mockService, mockOrgRepo, nil, nil)
		r := chi.NewRouter()
		r.Patch("/tickets/{ticketID}", h.UpdateTicket)

		ticketID := uuid.New()
		orgID := uuid.New()
		user := &domain.User{ID: uuid.New()}

		ticket := &domain.Ticket{ID: ticketID, OrganizationID: orgID}
		mockService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil)

		memberships := []domain.UserMembership{{Organization: domain.Organization{ID: orgID}}}
		mockOrgRepo.On("ListByUser", mock.Anything, user.ID).Return(memberships, nil)

		mockService.On("UpdateTicket", mock.Anything, ticketID, mock.Anything).Return(ticket, nil)

		reqBody := map[string]string{
			"title": "Valid Title",
			"description": "Valid Description",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PATCH", "/tickets/"+ticketID.String(), bytes.NewReader(bodyBytes))
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
