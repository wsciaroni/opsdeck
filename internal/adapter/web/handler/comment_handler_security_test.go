package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestCreateComment_Security(t *testing.T) {
	ticketID := uuid.New()
	orgID := uuid.New()
	user := &domain.User{ID: uuid.New()}
	ticket := &domain.Ticket{ID: ticketID, OrganizationID: orgID}

	t.Run("DoS Prevention - Rejects body too large", func(t *testing.T) {
		mockCommentService := new(MockCommentService)
		mockTicketService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		mockUserRepo := new(MockUserRepo)
		h := handler.NewCommentHandler(mockCommentService, mockTicketService, mockUserRepo, mockOrgRepo, nil)

		r := chi.NewRouter()
		r.Post("/tickets/{ticketID}/comments", h.Create)

		// Setup necessary mocks for flow up to body reading
		mockTicketService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil)
		mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil)

		// Create a reader larger than MaxJSONBodySize (1MB)
		// We use handler.MaxJSONBodySize as the expected limit.
		largeReader := &LargeReader{Size: handler.MaxJSONBodySize + 1024}

		req := httptest.NewRequest("POST", "/tickets/"+ticketID.String()+"/comments", largeReader)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	})

	t.Run("Validation - Rejects body too long", func(t *testing.T) {
		mockCommentService := new(MockCommentService)
		mockTicketService := new(MockTicketService)
		mockOrgRepo := new(MockOrgRepo)
		mockUserRepo := new(MockUserRepo)
		h := handler.NewCommentHandler(mockCommentService, mockTicketService, mockUserRepo, mockOrgRepo, nil)

		r := chi.NewRouter()
		r.Post("/tickets/{ticketID}/comments", h.Create)

		// Setup necessary mocks
		mockTicketService.On("GetTicket", mock.Anything, ticketID).Return(ticket, nil)
		mockOrgRepo.On("GetMemberRole", mock.Anything, orgID, user.ID).Return("member", nil)

		// Create body > 5000 chars
		longBody := LargeString(5001)
		reqBody := map[string]interface{}{
			"body":      longBody,
			"sensitive": false,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/tickets/"+ticketID.String()+"/comments", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "too long")
	})
}
